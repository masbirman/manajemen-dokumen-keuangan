<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import api from "@/services/api";
import { useAuthStore } from "@/stores/auth";

const authStore = useAuthStore();

const stats = ref({
  totalDokumen: 0,
  totalUnitKerja: 0,
  totalPPTK: 0,
  totalUser: 0,
});

// For operator - store assigned unit kerja and pptk names
const assignedUnitKerja = ref("");
const assignedPPTKList = ref<string[]>([]);

const recentDokumen = ref<
  Array<{
    id: string;
    nomor_dokumen: string;
    tanggal_dokumen: string;
    uraian: string;
    nilai: number;
  }>
>([]);

const roleLabels: Record<string, string> = {
  super_admin: "Super Admin",
  admin: "Admin",
  operator: "Operator",
};

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(
    value
  );

const formatDate = (date: string) => new Date(date).toLocaleDateString("id-ID");

const fetchStats = async () => {
  try {
    const [dokRes, ukRes, pptkRes, userRes] = await Promise.all([
      api.get("/dokumen").catch(() => ({ data: { data: [], total: 0 } })),
      authStore.isAdmin
        ? api.get("/unit-kerja").catch(() => ({ data: { data: [], total: 0 } }))
        : api.get("/unit-kerja/active").catch(() => ({ data: { data: [] } })),
      authStore.isAdmin
        ? api.get("/pptk").catch(() => ({ data: { data: [], total: 0 } }))
        : api.get("/pptk/active").catch(() => ({ data: { data: [] } })),
      authStore.isSuperAdmin
        ? api.get("/users").catch(() => ({ data: { data: [], total: 0 } }))
        : Promise.resolve({ data: { data: [], total: 0 } }),
    ]);

    // Get total from response or count array
    stats.value.totalDokumen =
      dokRes.data.total || (dokRes.data.data || dokRes.data || []).length;

    const unitKerjaList = ukRes.data.data || ukRes.data || [];
    const pptkList = pptkRes.data.data || pptkRes.data || [];

    if (authStore.isAdmin) {
      stats.value.totalUnitKerja = ukRes.data.total || unitKerjaList.length;
      stats.value.totalPPTK = pptkRes.data.total || pptkList.length;
    } else {
      // For operator, find assigned unit kerja and pptk names
      if (authStore.user?.unit_kerja_id) {
        const uk = unitKerjaList.find(
          (u: { id: string; nama: string }) =>
            u.id === authStore.user?.unit_kerja_id
        );
        assignedUnitKerja.value = uk?.nama || "-";
      }
      // Get PPTK names from pptk_list
      if (authStore.user?.pptk_list?.length) {
        const pptkNames = authStore.user.pptk_list
          .map((up: { pptk_id: string; pptk?: { nama: string } }) => {
            if (up.pptk?.nama) return up.pptk.nama;
            const found = pptkList.find(
              (p: { id: string; nama: string }) => p.id === up.pptk_id
            );
            return found?.nama;
          })
          .filter(Boolean);
        assignedPPTKList.value = pptkNames;
      } else if (authStore.user?.pptk_id) {
        // Fallback to single pptk_id
        const pptk = pptkList.find(
          (p: { id: string; nama: string }) => p.id === authStore.user?.pptk_id
        );
        if (pptk?.nama) assignedPPTKList.value = [pptk.nama];
      }
    }

    stats.value.totalUser =
      userRes.data.total || (userRes.data.data || userRes.data || []).length;

    const dokumen = dokRes.data.data || dokRes.data || [];
    recentDokumen.value = dokumen.slice(0, 5);
  } catch {
    // ignore
  }
};

onMounted(fetchStats);
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-gray-800">Dashboard</h1>
      <p class="text-gray-500">Selamat datang, {{ authStore.user?.name }}</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Total Dokumen</p>
            <p class="text-2xl font-bold text-gray-800">
              {{ stats.totalDokumen }}
            </p>
          </div>
          <div
            class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center text-2xl"
          >
            📄
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Unit Kerja</p>
            <p
              v-if="authStore.isAdmin"
              class="text-2xl font-bold text-gray-800"
            >
              {{ stats.totalUnitKerja }}
            </p>
            <p
              v-else
              class="text-lg font-bold text-gray-800 truncate max-w-[150px]"
              :title="assignedUnitKerja"
            >
              {{ assignedUnitKerja || "-" }}
            </p>
          </div>
          <div
            class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center text-2xl"
          >
            🏢
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div class="flex items-center justify-between">
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-500">PPTK</p>
            <p
              v-if="authStore.isAdmin"
              class="text-2xl font-bold text-gray-800"
            >
              {{ stats.totalPPTK }}
            </p>
            <div v-else>
              <p
                v-if="assignedPPTKList.length === 1"
                class="text-lg font-bold text-gray-800 truncate"
                :title="assignedPPTKList[0]"
              >
                {{ assignedPPTKList[0] }}
              </p>
              <div v-else-if="assignedPPTKList.length > 1" class="space-y-1">
                <p
                  v-for="(name, idx) in assignedPPTKList.slice(0, 2)"
                  :key="idx"
                  class="text-sm font-medium text-gray-800 truncate"
                >
                  {{ name }}
                </p>
                <p
                  v-if="assignedPPTKList.length > 2"
                  class="text-xs text-gray-500"
                  :title="assignedPPTKList.join(', ')"
                >
                  +{{ assignedPPTKList.length - 2 }} lagi
                </p>
              </div>
              <p v-else class="text-lg font-bold text-gray-800">-</p>
            </div>
          </div>
          <div
            class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center text-2xl flex-shrink-0"
          >
            👤
          </div>
        </div>
      </div>

      <div
        v-if="authStore.isSuperAdmin"
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-100"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Total User</p>
            <p class="text-2xl font-bold text-gray-800">
              {{ stats.totalUser }}
            </p>
          </div>
          <div
            class="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center text-2xl"
          >
            👥
          </div>
        </div>
      </div>

      <div
        v-else
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-100"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">Role Anda</p>
            <p class="text-xl font-bold text-gray-800">
              {{ roleLabels[authStore.user?.role || ""] }}
            </p>
          </div>
          <div
            class="w-12 h-12 bg-yellow-100 rounded-full flex items-center justify-center text-2xl"
          >
            🔑
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">Aksi Cepat</h2>
        <div class="grid grid-cols-2 gap-3">
          <RouterLink
            to="/dokumen/input"
            class="p-4 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors text-center"
          >
            <span class="text-2xl block mb-1">📝</span>
            <span class="text-sm font-medium text-blue-700">Input Dokumen</span>
          </RouterLink>
          <RouterLink
            to="/dokumen"
            class="p-4 bg-green-50 rounded-lg hover:bg-green-100 transition-colors text-center"
          >
            <span class="text-2xl block mb-1">📋</span>
            <span class="text-sm font-medium text-green-700">List Dokumen</span>
          </RouterLink>
          <RouterLink
            v-if="authStore.isAdmin"
            to="/unit-kerja"
            class="p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors text-center"
          >
            <span class="text-2xl block mb-1">🏢</span>
            <span class="text-sm font-medium text-purple-700">Unit Kerja</span>
          </RouterLink>
          <RouterLink
            v-if="authStore.isAdmin"
            to="/pptk"
            class="p-4 bg-yellow-50 rounded-lg hover:bg-yellow-100 transition-colors text-center"
          >
            <span class="text-2xl block mb-1">👤</span>
            <span class="text-sm font-medium text-yellow-700">PPTK</span>
          </RouterLink>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">
          Dokumen Terbaru
        </h2>
        <div
          v-if="recentDokumen.length === 0"
          class="text-gray-500 text-center py-4"
        >
          Belum ada dokumen
        </div>
        <div v-else class="space-y-3">
          <RouterLink
            v-for="doc in recentDokumen"
            :key="doc.id"
            to="/dokumen"
            class="block p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="flex justify-between items-start">
              <div>
                <p class="font-medium text-gray-800">{{ doc.nomor_dokumen }}</p>
                <p class="text-sm text-gray-500 truncate max-w-xs">
                  {{ doc.uraian }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium text-blue-600">
                  {{ formatCurrency(doc.nilai) }}
                </p>
                <p class="text-xs text-gray-400">
                  {{ formatDate(doc.tanggal_dokumen) }}
                </p>
              </div>
            </div>
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>
