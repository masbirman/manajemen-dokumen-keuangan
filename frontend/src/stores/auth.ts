import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/services/api";

export interface UserPPTK {
  id: string;
  pptk_id: string;
  pptk?: {
    id: string;
    nama: string;
  };
}

export interface User {
  id: string;
  username: string;
  name: string;
  role: "super_admin" | "admin" | "operator";
  unit_kerja_id?: string;
  pptk_id?: string;
  pptk_list?: UserPPTK[];
  avatar_path?: string;
  is_active: boolean;
}

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem("access_token"));
  const loading = ref(false);

  const isAuthenticated = computed(() => !!token.value && !!user.value);
  const isAdmin = computed(
    () => user.value?.role === "admin" || user.value?.role === "super_admin"
  );
  const isSuperAdmin = computed(() => user.value?.role === "super_admin");
  const isOperator = computed(() => user.value?.role === "operator");

  async function login(username: string, password: string) {
    loading.value = true;
    try {
      const response = await api.post("/auth/login", { username, password });
      const tokenData = response.data.data?.token || response.data;
      const accessToken = tokenData.access_token || tokenData.token;

      token.value = accessToken;
      localStorage.setItem("access_token", accessToken);
      if (tokenData.refresh_token) {
        localStorage.setItem("refresh_token", tokenData.refresh_token);
      }

      // Set user from response if available
      if (response.data.data?.user) {
        user.value = response.data.data.user;
      } else {
        await fetchUser();
      }
    } finally {
      loading.value = false;
    }
  }

  async function logout() {
    try {
      await api.post("/auth/logout");
    } finally {
      token.value = null;
      user.value = null;
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
    }
  }

  async function fetchUser() {
    if (!token.value) return;
    try {
      const response = await api.get("/auth/me");
      user.value = response.data.data || response.data;
    } catch {
      token.value = null;
      user.value = null;
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
    }
  }

  async function initialize() {
    if (token.value && !user.value) {
      await fetchUser();
    }
  }

  function hasRole(roles: string[]): boolean {
    return user.value ? roles.includes(user.value.role) : false;
  }

  return {
    user,
    token,
    loading,
    isAuthenticated,
    isAdmin,
    isSuperAdmin,
    isOperator,
    login,
    logout,
    fetchUser,
    initialize,
    hasRole,
  };
});
