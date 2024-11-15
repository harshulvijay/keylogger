import { Router } from "express";
import log from "./log.ts";
import ws from "./ws.ts";

// router
const router = Router({});

router.use("/log", log);
router.use("/ws", ws);

export default router;
