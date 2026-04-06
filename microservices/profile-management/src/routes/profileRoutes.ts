import { Router, Request, Response } from "express";
import { AppDataSource } from "../config/db";
import { Profile } from "../models/Profile";
import { z } from "zod";

const router = Router();

// Local middleware removed - security handled by API Gateway

const UpdateProfileSchema = z.object({
  name: z.string().optional(),
  email: z.string().email().optional(),
  phone: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  bio: z.string().optional(),
});

// GET /profiles - Returns all the profiles paginated.
router.get("/", async (req: Request, res: Response) => {
  try {
    const page = parseInt(req.query.page as string) || 1;
    const limit = parseInt(req.query.limit as string) || 10;
    const skip = (page - 1) * limit;

    const profileRepository = AppDataSource.getRepository(Profile);
    const [profiles, total] = await profileRepository.findAndCount({
      skip,
      take: limit,
      order: { createdAt: "DESC" },
    });

    res.json({
      data: profiles,
      pagination: {
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      },
    });
  } catch (error) {
    console.error("Error fetching profiles:", error);
    res.status(500).json({ message: "Internal server error" });
  }
});

// GET /profiles/{employeeId} - Consult the requested employee.
router.get("/:employeeId", async (req: Request, res: Response) => {
  try {
    const { employeeId } = req.params;
    const profileRepository = AppDataSource.getRepository(Profile);
    const profile = await profileRepository.findOne({
      where: { employeeId },
    });

    if (!profile) {
      return res.status(404).json({ message: `Profile for employee ${employeeId} not found` });
    }

    res.json(profile);
  } catch (error) {
    console.error("Error fetching profile:", error);
    res.status(500).json({ message: "Internal server error" });
  }
});

// PUT /profiles/{employeeId} - Update fully or partially update an employees updatable information.
router.put("/:employeeId", async (req: Request, res: Response) => {
  try {
    const { employeeId } = req.params;
    const validation = UpdateProfileSchema.safeParse(req.body);

    if (!validation.success) {
      return res.status(400).json({
        message: "Invalid input data",
        errors: validation.error.errors,
      });
    }

    const profileRepository = AppDataSource.getRepository(Profile);
    const profile = await profileRepository.findOne({
      where: { employeeId },
    });

    if (!profile) {
      return res.status(404).json({ message: `Profile for employee ${employeeId} not found` });
    }

    // Update profile
    Object.assign(profile, validation.data);
    await profileRepository.save(profile);

    res.json(profile);
  } catch (error) {
    console.error("Error updating profile:", error);
    res.status(500).json({ message: "Internal server error" });
  }
});

export default router;
