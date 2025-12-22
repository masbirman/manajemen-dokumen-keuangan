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
  const selectedYear = ref<string>(localStorage.getItem("selectedYear") || "2025");
  const loading = ref(false);

  const isAuthenticated = computed(() => !!token.value && !!user.value);
  const isAdmin = computed(
    () => user.value?.role === "admin" || user.value?.role === "super_admin"
  );
  const isSuperAdmin = computed(() => user.value?.role === "super_admin");
  const isOperator = computed(() => user.value?.role === "operator");

  async function login(username: string, password: string, tahunAnggaran = '2025', turnstileToken = '') {
    loading.value = true;
    try {
      const payload: any = { username, password, tahun_anggaran: tahunAnggaran };
      if (turnstileToken) {
        payload.turnstile_token = turnstileToken;
      }

      const response = await api.post("/auth/login", payload);
      const tokenData = response.data.data?.token || response.data;
      const accessToken = tokenData.access_token || tokenData.token;

      token.value = accessToken;
      user.value = response.data.data?.user || response.data.user;
      selectedYear.value = tahunAnggaran;

      localStorage.setItem("access_token", accessToken);
      localStorage.setItem("selectedYear", tahunAnggaran);
      
      if (tokenData.refresh_token) {
        localStorage.setItem("refresh_token", tokenData.refresh_token);
      }
      
      if (!user.value) {
         await fetchUser();
      }

      return { success: true };
    } catch (error: any) {
      return { 
        success: false, 
        message: error.response?.data?.message || 'Login failed' 
      };
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
      localStorage.removeItem("selectedYear");
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

  async function checkAuth() {
    await fetchUser();
  }

  async function initialize() {
    if (token.value && !user.value) {
      await fetchUser();
    }
  }

  function hasRole(roles: string[] | string): boolean {
    if (Array.isArray(roles)) {
        return user.value ? roles.includes(user.value.role) : false;
    }
    return user.value?.role === roles;
  }

  function hasAnyRole(roles: string[]): boolean {
    return hasRole(roles);
  }

  function setYear(year: string) {
    selectedYear.value = year;
    localStorage.setItem("selectedYear", year);
  }

  return {
    user,
    token,
    selectedYear,
    loading,
    isAuthenticated,
    isAdmin,
    isSuperAdmin,
    isOperator,
    login,
    logout,
    fetchUser,
    checkAuth,
    initialize,
    hasRole,
    hasAnyRole,
    setYear
  };
});
