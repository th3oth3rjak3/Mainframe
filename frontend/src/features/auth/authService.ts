import * as v from "valibot";
import { httpClient } from "@/lib/httpClient";
import { handleApiRequest } from "@/lib/apiHelpers";
import { type LoginResponse, LoginResponseSchema, type SignInRequest } from "@/features/auth/types";

// @ts-ignore
import type { ApiError } from "@/lib/apiHelpers";
// @ts-ignore
import type { ValiError } from "valibot";

export interface IAuthService {
  /**
   * Sign in a user with username and password
   * @throws {ApiError} When authentication fails or server returns an error
   * @throws {ValiError} When the server response doesn't match the expected schema.
   */
  signIn(request: SignInRequest): Promise<LoginResponse>;

  /**
   * Sign out the current user
   * @throws {ApiError} When the request fails
   */
  signOut(): Promise<void>;
}

export const authService: IAuthService = {
  signIn: async (request) => {
    return handleApiRequest(async () => {
      const response = await httpClient.post("auth/login", { json: request }).json();

      return v.parse(LoginResponseSchema, response);
    });
  },

  signOut: async () => {
    return handleApiRequest(async () => {
      await httpClient.post("auth/logout");
    });
  },
};
