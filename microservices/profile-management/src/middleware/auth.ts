import { Request, Response, NextFunction } from "express";
import * as crypto from "crypto";

const JWT_SECRET = process.env.JWT_SECRET || "supersecretkey";

export const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
  const authHeader = req.headers.authorization;

  if (!authHeader) {
    return res.status(401).json({ message: "Authorization header is required" });
  }

  const [bearer, token] = authHeader.split(" ");

  if (bearer !== "Bearer" || !token) {
    return res.status(401).json({ message: "Invalid authorization header format" });
  }

  try {
    const [headerB64, payloadB64, signatureB64] = token.split(".");
    
    if (!headerB64 || !payloadB64 || !signatureB64) {
        return res.status(401).json({ message: "Invalid token format" });
    }

    // Verify signature
    const expectedSignature = crypto
      .createHmac("sha256", JWT_SECRET)
      .update(`${headerB64}.${payloadB64}`)
      .digest("base64url")
      .replace(/=/g, "");

    // signatureB64 from JWT might be base64url or base64.
    // Standard JWT uses base64url.
    if (signatureB64 !== expectedSignature) {
      return res.status(401).json({ message: "Invalid token signature" });
    }

    const payload = JSON.parse(Buffer.from(payloadB64, "base64").toString());
    
    // Check expiration
    if (payload.exp && Date.now() >= payload.exp * 1000) {
      return res.status(401).json({ message: "Token has expired" });
    }

    const role = payload.role;
    
    // RBAC: ADMIN has total access. USER has read-only access.
    if (req.method !== "GET") {
      if (role !== "ADMIN") {
        return res.status(403).json({ message: "Forbidden: ADMIN role required for write operations" });
      }
    }

    // Attach user info to request
    (req as any).user = payload;
    
    next();
  } catch (error) {
    console.error("JWT validation error:", error);
    return res.status(401).json({ message: "Invalid or expired token" });
  }
};
