<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from "vue";
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
    jenis_dokumen?: { nama: string; kode?: string };
    sumber_dana?: { nama: string; kode?: string };
  }>
>([]);

const roleLabels: Record<string, string> = {
  super_admin: "Super Admin",
  admin: "Admin Keuangan",
  operator: "Operator",
};

// Dynamic greeting based on time
const getGreeting = () => {
  const hour = new Date().getHours();
  if (hour >= 5 && hour < 11) return "Selamat Pagi";
  if (hour >= 11 && hour < 15) return "Selamat Siang";
  if (hour >= 15 && hour < 18) return "Selamat Sore";
  return "Selamat Malam";
};

const greeting = ref(getGreeting());

// Typing effect
const displayText = ref("");
const isTyping = ref(true);
const typingIntervalId = ref<number | null>(null);

const fullGreetingText = computed(() => {
  return `${greeting.value}, ${authStore.user?.name || 'User'}`;
});

const startTypingEffect = () => {
  let charIndex = 0;
  const text = fullGreetingText.value;
  displayText.value = "";
  isTyping.value = true;
  
  // Clear previous interval
  if (typingIntervalId.value) {
    clearInterval(typingIntervalId.value);
  }
  
  typingIntervalId.value = window.setInterval(() => {
    if (charIndex < text.length) {
      displayText.value += text.charAt(charIndex);
      charIndex++;
    } else {
      isTyping.value = false;
      // Wait 3 seconds then restart
      setTimeout(() => {
        // Update greeting in case time changed
        greeting.value = getGreeting();
        startTypingEffect();
      }, 3000);
      if (typingIntervalId.value) {
        clearInterval(typingIntervalId.value);
      }
    }
  }, 80);
};

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR" }).format(
    value
  );

const formatDate = (date: string) => new Date(date).toLocaleDateString("id-ID");

const logoUrl = ref("");

// Fetch logo settings
const fetchLogo = async () => {
  try {
    const response = await api.get("/public/login-settings");
    const data = response.data.data || {};
    if (data.login_logo_url) {
      const url = data.login_logo_url;
      if (url.startsWith("/")) {
        const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8000/api";
        const baseUrl = apiUrl.replace(/\/api$/, "");
        logoUrl.value = baseUrl + url;
      } else {
        logoUrl.value = url;
      }
    }
  } catch {
    // ignore
  }
};

const fetchStats = async () => {
  try {
    // Ensure we have the latest user data including relations
    await authStore.fetchUser();

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

onMounted(() => {
  fetchStats();
  fetchLogo();
  startTypingEffect();
});

onUnmounted(() => {
  if (typingIntervalId.value) {
    clearInterval(typingIntervalId.value);
  }
});
</script>

<template>
  <div>
    <!-- Hero Welcome Section -->
    <!-- Hero Welcome Section -->
    <div class="mb-8 bg-indigo-600 rounded-2xl p-6 md:p-8 shadow-md relative overflow-hidden">
      
      <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div>
          <!-- Typing greeting -->
          <h1 class="text-2xl md:text-3xl font-bold text-white mb-3">
            {{ displayText }}<span class="animate-pulse" :class="isTyping ? 'opacity-100' : 'opacity-0'">|</span>
          </h1>
          
          <!-- Info badges -->
          <div class="flex flex-wrap gap-2 md:gap-3 mt-4">
            <span class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-700/50 border border-indigo-500/30 rounded-full text-white text-sm font-medium">
              <span class="w-1.5 h-1.5 bg-green-400 rounded-full animate-pulse"></span>
              {{ roleLabels[authStore.user?.role || ""] }}
            </span>
            <span v-if="assignedUnitKerja || authStore.isAdmin" class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-700/50 border border-indigo-500/30 rounded-full text-white text-sm font-medium">
              🏢 {{ authStore.isAdmin ? `${stats.totalUnitKerja} Unit Kerja` : assignedUnitKerja }}
            </span>
            <span v-if="assignedPPTKList.length > 0 || authStore.isAdmin" class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-indigo-700/50 border border-indigo-500/30 rounded-full text-white text-sm font-medium">
              👤 {{ authStore.isAdmin ? `${stats.totalPPTK} PPTK` : assignedPPTKList.join(', ') }}
            </span>
          </div>
        </div>

        <!-- Logo or Clock icon -->
        <div class="flex-shrink-0">
          <img
            v-if="logoUrl"
            :src="logoUrl"
            alt="Logo"
            class="h-16 md:h-20 w-auto object-contain drop-shadow-md"
          />
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            class="w-20 h-20 text-indigo-400 opacity-50"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="1.5"
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow">
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

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow">
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

      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow">
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
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow"
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
        class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow"
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
                <p class="font-medium text-gray-800">
                  {{ doc.jenis_dokumen?.kode || doc.jenis_dokumen?.nama || 'Jenis' }} / {{ doc.sumber_dana?.nama || 'Sumber Dana' }}
                </p>
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
