import { Router } from "express";
import log from "./log";
import ws from "./ws";

// router
const router = Router({});

router.use("/log", log);
router.use("/ws", ws);

export default router;
