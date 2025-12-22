import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
    meta: { guest: true },
  },
  {
    path: "/",
    component: () => import("@/layouts/MainLayout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => import("@/views/DashboardView.vue"),
      },
      {
        path: "unit-kerja",
        name: "unit-kerja",
        component: () => import("@/views/UnitKerjaView.vue"),
        meta: { roles: ["super_admin", "admin"] },
      },
      {
        path: '/profile',
        name: 'profile',
        component: () => import('@/views/ProfileView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: '/:pathMatch(.*)*',
        redirect: '/'
      },
      {
        path: "pptk",
        name: "pptk",
        component: () => import("@/views/PPTKView.vue"),
        meta: { roles: ["super_admin", "admin"] },
      },
      {
        path: "sumber-dana",
        name: "sumber-dana",
        component: () => import("@/views/SumberDanaView.vue"),
        meta: { roles: ["super_admin", "admin"] },
      },
      {
        path: "jenis-dokumen",
        name: "jenis-dokumen",
        component: () => import("@/views/JenisDokumenView.vue"),
        meta: { roles: ["super_admin", "admin"] },
      },
      {
        path: "users",
        name: "users",
        component: () => import("@/views/ManajemenUserView.vue"),
        meta: { roles: ["super_admin"] },
      },
      {
        path: "dokumen",
        name: "dokumen-list",
        component: () => import("@/views/ListDokumenView.vue"),
      },
      {
        path: "dokumen/input",
        name: "dokumen-input",
        component: () => import("@/views/InputDokumenView.vue"),
        meta: { roles: ["super_admin", "admin", "operator"] },
      },
      {
        path: "dokumen/edit/:id",
        name: "dokumen-edit",
        component: () => import("@/views/EditDokumenView.vue"),
        meta: { roles: ["super_admin", "admin", "operator"] },
      },
      {
        path: "settings",
        name: "settings",
        component: () => import("@/views/PengaturanView.vue"),
        meta: { roles: ["super_admin"] },
      },
      {
        path: "petunjuk",
        name: "petunjuk",
        component: () => import("@/views/PetunjukView.vue"),
        meta: { roles: ["super_admin"] },
      },
      {
        path: "login-settings",
        name: "login-settings",
        component: () => import("@/views/PengaturanLoginView.vue"),
        meta: { roles: ["super_admin"] },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();

  // Initialize auth state if needed
  await authStore.initialize();

  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const isGuestOnly = to.matched.some((record) => record.meta.guest);
  const requiredRoles = to.meta.roles as string[] | undefined;

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: "login" });
    return;
  }

  if (isGuestOnly && authStore.isAuthenticated) {
    next({ name: "dashboard" });
    return;
  }

  if (requiredRoles && !authStore.hasRole(requiredRoles)) {
    next({ name: "dashboard" });
    return;
  }

  next();
});

export default router;
