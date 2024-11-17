import { randomUUID } from "crypto";
import { Router } from "express";
import formidable from "formidable";
import { readFile } from "fs/promises";
import GITHUB_TOKEN from "../util/octokit.ts";

// router
const router = Router({});

const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/gif"];

router.route("/save").post(async (req, res) => {
	const url = `https://api.github.com/repos/${process.env.GITHUB_OWNER}/${process.env.GITHUB_REPOSITORY}/contents/`;

	const form = formidable({
		allowEmptyFiles: false,
		multiples: false,
		keepExtensions: true,
	});

	form.parse(req, async (err, fields, files) => {
		if (err) {
			res.status(500).json({
				errors: ["Cannot parse form data", err],
			});
			return;
		}

		if (files) {
			// make sure the files are sent with the `files` key in form data
			if (!files.files) {
				// bad request
				res.status(400);
				res.end();
				return;
			}

			const targetID = fields["t"]?.[0] || "unknown";
			const branch = fields["b"]?.[0] || "main";

			try {
				files.files.map(async (file) => {
					const fileName = `${new Date().getTime()} ${
						file.originalFilename || randomUUID()
					}`;
					const mimetype = file.mimetype;

					// strictly allow image files only
					if (!ALLOWED_IMAGE_TYPES.includes(mimetype || "")) {
						// bad request
						res.status(400);
						res.end();
						return;
					}

					const content = file.filepath
						? await readFile(file.filepath, "binary")
						: file.toString();
					const encodedContent = Buffer.from(content, "binary").toString(
						"base64"
					);

					const response = await fetch(`${url}${fileName}`, {
						method: "PUT",
						headers: {
							"Authorization": `Bearer ${GITHUB_TOKEN}`,
							"Content-Type": "application/json",
						},
						body: JSON.stringify({
							message: `screenshot modified at: ${
								file.mtime || "<unknown>"
							}\ntarget ID: ${targetID}`,
							content: encodedContent,
							branch,
						}),
					});

					res.status(response.status);
					res.end();
					return;
				});
			} catch {
				// internal server error
				res.status(500);
				res.end();
				return;
			}
		}
	});
});

export default router;
