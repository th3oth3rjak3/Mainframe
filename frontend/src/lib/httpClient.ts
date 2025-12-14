import ky from "ky";

/**
 * httpClient is used to make HTTP requests to the backend API.
 * All requests are automatically prefixed with '/api' and include credentials (cookies).
 *
 * @example
 * // POST to /api/auth/signout
 * await httpClient.post('auth/signout')
 *
 * @example
 * // GET from /api/users/me
 * const user = await httpClient.get('users/me').json()
 */
export const httpClient = ky.create({
  prefixUrl: "/api",
  credentials: "include",
  timeout: 10_000,
});
