// src/lib/apiHelpers.ts
import { HTTPError } from "ky";
import * as v from "valibot";

/**
 * Error response schema matching the backend ErrorResponse struct
 */
export const ErrorResponseSchema = v.object({
  error: v.string(),
});

export type ErrorResponse = v.InferOutput<typeof ErrorResponseSchema>;

/**
 * Custom error class for API-related errors
 */
export class ApiError extends Error {
  statusCode: number;
  response: ErrorResponse;

  constructor(statusCode: number, response: ErrorResponse) {
    super(response.error);
    this.name = "ApiError";
    this.statusCode = statusCode;
    this.response = response;
  }
}

/**
 * Wraps an HTTP request and transforms errors into ApiError instances
 * @param fn - The async function that makes the HTTP request
 * @returns The result of the function
 * @throws {ApiError} When the HTTP request fails
 */
export async function handleApiRequest<T>(fn: () => Promise<T>): Promise<T> {
  try {
    return await fn();
  } catch (error) {
    console.log(error);
    if (error instanceof HTTPError) {
      const errorBody = await error.response.json();
      const errorResponse = v.parse(ErrorResponseSchema, errorBody);
      throw new ApiError(error.response.status, errorResponse);
    }
    throw error;
  }
}
