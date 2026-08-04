import { reactive } from "vue";

export const authState = reactive({
  username: sessionStorage.getItem("username") || null,
  userId: sessionStorage.getItem("user_id") ? Number(sessionStorage.getItem("user_id")) : null,
  isAuthenticated: !!sessionStorage.getItem("access_token"),
});

export function setSession(username, userId, accessToken, refreshToken) {
  sessionStorage.setItem("username", username);
  sessionStorage.setItem("user_id", String(userId));
  sessionStorage.setItem("access_token", accessToken);
  sessionStorage.setItem("refresh_token", refreshToken);
  authState.username = username;
  authState.userId = userId;
  authState.isAuthenticated = true;
}

export function clearSession() {
  sessionStorage.removeItem("username");
  sessionStorage.removeItem("user_id");
  sessionStorage.removeItem("access_token");
  sessionStorage.removeItem("refresh_token");
  authState.username = null;
  authState.userId = null;
  authState.isAuthenticated = false;
}
