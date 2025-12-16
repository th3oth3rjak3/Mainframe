import { create } from "zustand";
import { authService, type IAuthService } from "./authService";
import { AuthenticatedUser, type SignInRequest } from "./types";
// @ts-ignore
import { ApiError } from "@/lib/apiHelpers";

type AuthStore = {
  user: AuthenticatedUser | null;
  isLoading: boolean;
  error: string | null;

  /**
   * Sign in a user with username and password
   * @throws {ApiError} When authentication fails
   */
  signIn: (request: SignInRequest, service?: IAuthService) => Promise<void>;

  /**
   * Sign out the current user
   */
  signOut: (service?: IAuthService) => Promise<void>;

  /**
   * Check if a user is authenticated on mount
   */
  checkAuth: (service?: IAuthService) => Promise<void>;
};

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  isLoading: true,
  error: null,

  signIn: async (request, service = authService) => {
    try {
      set({ error: null });
      const response = await service.signIn(request);
      const user = new AuthenticatedUser(response);
      set({ user, error: null });
    } catch (error) {
      const message = error instanceof ApiError ? error.message : "Sign in failed";
      set({ error: message, user: null });
      throw error;
    }
  },

  signOut: async (service = authService) => {
    try {
      await service.signOut();
      set({ user: null, error: null });
    } catch (error) {
      const message = error instanceof ApiError ? error.message : "Sign out failed";
      set({ error: message });
      throw error;
    }
  },

  checkAuth: async (service = authService) => {
    // TODO: Add a getCurrentUser method to authService when backend endpoint exists
    // @ts-ignore
    const _ = service;
    set({ isLoading: false });
  },
}));
