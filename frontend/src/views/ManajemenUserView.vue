<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { useAuthStore } from "@/stores/auth";
import { DataTable, Modal, InputField, Dropdown } from "@/components/ui";

interface UnitKerja {
  id: string;
  nama: string;
}
interface PPTK {
  id: string;
  nama: string;
  unit_kerja_id: string;
}
interface UserPPTK {
  id: string;
  pptk_id: string;
  pptk?: PPTK;
}
interface User {
  id: string;
  username: string;
  name: string;
  role: string;
  is_active: boolean;
  unit_kerja_id?: string;
  pptk_id?: string;
  unit_kerja?: UnitKerja;
  pptk?: PPTK;
  pptk_list?: UserPPTK[];
  avatar_path?: string;
}

const toast = useToast();
const authStore = useAuthStore();
const loading = ref(false);
const data = ref<User[]>([]);
const unitKerjaList = ref<UnitKerja[]>([]);
const pptkList = ref<PPTK[]>([]);
const showModal = ref(false);
const showDeleteModal = ref(false);
const showBulkDeleteModal = ref(false);
const showAvatarModal = ref(false);
const editingItem = ref<User | null>(null);
const deletingItem = ref<User | null>(null);
const avatarUser = ref<User | null>(null);
const selectedIds = ref<string[]>([]);
const avatarFile = ref<File | null>(null);
const avatarPreview = ref<string>("");
const showResetPasswordModal = ref(false);
const resetPasswordResult = ref<{ username: string; new_password: string } | null>(null);
const resettingPassword = ref(false);

// Pagination & Search
const currentPage = ref(1);
const totalPages = ref(1);
const totalItems = ref(0);
const perPage = ref(10);
const searchQuery = ref("");
const filterRole = ref("");

const form = ref({
  username: "",
  name: "",
  password: "",
  role: "operator",
  unit_kerja_id: null as string | null,
  pptk_id: null as string | null,
  pptk_ids: [] as string[],
});
const errors = ref<Record<string, string>>({});

const columns = [
  { key: "avatar", label: "", width: "60px" },
  { key: "username", label: "Username", width: "150px" },
  { key: "name", label: "Nama" },
  { key: "role", label: "Role", width: "100px" },
  { key: "unit_kerja", label: "Unit Kerja" },
  { key: "pptk_list", label: "PPTK" },
  { key: "status", label: "Status", width: "80px" },
  { key: "actions", label: "Aksi", width: "180px" },
];

const roleOptions = [
  { value: "super_admin", label: "Super Admin" },
  { value: "admin", label: "Admin" },
  { value: "operator", label: "Operator" },
];

const filterRoleOptions = [
  { value: "", label: "Semua Role" },
  { value: "super_admin", label: "Super Admin" },
  { value: "admin", label: "Admin" },
  { value: "operator", label: "Operator" },
];

const unitKerjaOptions = computed(() =>
  unitKerjaList.value.map((uk) => ({ value: uk.id, label: uk.nama }))
);
const filteredPptkOptions = computed(() =>
  pptkList.value
    .filter(
      (p) =>
        !form.value.unit_kerja_id ||
        p.unit_kerja_id === form.value.unit_kerja_id
    )
    .map((p) => ({ value: p.id, label: p.nama }))
);

const roleLabels: Record<string, string> = {
  super_admin: "Super Admin",
  admin: "Admin",
  operator: "Operator",
};

watch(
  () => form.value.unit_kerja_id,
  () => {
    form.value.pptk_id = null;
  }
);

const getAvatarUrl = (path?: string) => {
  if (!path) return "";
  return `${
    import.meta.env.VITE_API_URL || "http://localhost:8000/api"
  }/files/${path}`;
};

const fetchData = async (page = 1) => {
  loading.value = true;
  try {
    const params: Record<string, unknown> = { page, page_size: perPage.value };
    if (searchQuery.value) params.search = searchQuery.value;
    if (filterRole.value) params.role = filterRole.value;

    const [userRes, ukRes, pptkRes] = await Promise.all([
      api.get("/users", { params }),
      api.get("/unit-kerja/active"),
      api.get("/pptk/active"),
    ]);
    data.value = userRes.data.data || userRes.data || [];
    currentPage.value = userRes.data.page || 1;
    totalPages.value = userRes.data.total_pages || 1;
    totalItems.value = userRes.data.total || 0;
    unitKerjaList.value = ukRes.data.data || ukRes.data || [];
    pptkList.value = pptkRes.data.data || pptkRes.data || [];
  } catch {
    toast.error("Gagal memuat data");
  } finally {
    loading.value = false;
  }
};

const onPageChange = (page: number) => fetchData(page);

const onPerPageChange = (newPerPage: number) => {
  perPage.value = newPerPage;
  fetchData(1);
};

const onSearch = (query: string) => {
  searchQuery.value = query;
  fetchData(1);
};

const onFilterRole = () => {
  fetchData(1);
};

const openCreate = () => {
  editingItem.value = null;
  form.value = {
    username: "",
    name: "",
    password: "",
    role: "operator",
    unit_kerja_id: null,
    pptk_id: null,
    pptk_ids: [],
  };
  errors.value = {};
  showModal.value = true;
};

const openEdit = (item: User) => {
  editingItem.value = item;
  // Extract PPTK IDs from pptk_list
  const pptkIds = item.pptk_list?.map((up) => up.pptk_id) || [];
  form.value = {
    username: item.username,
    name: item.name,
    password: "",
    role: item.role,
    unit_kerja_id: item.unit_kerja_id || null,
    pptk_id: item.pptk_id || null,
    pptk_ids: pptkIds,
  };
  errors.value = {};
  showModal.value = true;
};

const openAvatarUpload = (user: User) => {
  avatarUser.value = user;
  avatarFile.value = null;
  avatarPreview.value = user.avatar_path ? getAvatarUrl(user.avatar_path) : "";
  showAvatarModal.value = true;
};

const onAvatarChange = (e: Event) => {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) {
    avatarFile.value = file;
    avatarPreview.value = URL.createObjectURL(file);
  }
};

const uploadAvatar = async () => {
  if (!avatarUser.value || !avatarFile.value) return;
  const formData = new FormData();
  formData.append("avatar", avatarFile.value);
  try {
    await api.post(`/users/${avatarUser.value.id}/avatar`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    toast.success("Avatar berhasil diupload");
    showAvatarModal.value = false;

    // Refresh auth store if uploading avatar for current user
    if (avatarUser.value.id === authStore.user?.id) {
      await authStore.fetchUser();
    }

    fetchData();
  } catch {
    toast.error("Gagal upload avatar");
  }
};

const save = async () => {
  errors.value = {};
  const payload: Record<string, unknown> = { ...form.value };
  if (editingItem.value && !payload.password) delete payload.password;
  try {
    if (editingItem.value) {
      await api.put(`/users/${editingItem.value.id}`, payload);
      toast.success("User berhasil diupdate");
    } else {
      await api.post("/users", payload);
      toast.success("User berhasil ditambahkan");
    }
    showModal.value = false;
    fetchData();
  } catch (e: unknown) {
    const err = e as {
      response?: { data?: { errors?: Record<string, string> } };
    };
    errors.value = err.response?.data?.errors || {};
    toast.error("Gagal menyimpan data");
  }
};

const toggleActive = async (user: User) => {
  try {
    if (user.is_active) {
      await api.delete(`/users/${user.id}`);
      toast.success("User dinonaktifkan");
    } else {
      await api.put(`/users/${user.id}`, { ...user, is_active: true });
      toast.success("User diaktifkan");
    }
    fetchData();
  } catch {
    toast.error("Gagal mengubah status");
  }
};

const confirmDelete = async () => {
  if (!deletingItem.value) return;
  try {
    await api.delete(`/users/${deletingItem.value.id}`);
    toast.success("User dinonaktifkan");
    showDeleteModal.value = false;
    fetchData();
  } catch {
    toast.error("Gagal menonaktifkan user");
  }
};

const onSelectionChange = (ids: (string | number)[]) => {
  selectedIds.value = ids as string[];
};

const confirmBulkDeactivate = async () => {
  try {
    await Promise.all(
      selectedIds.value.map((id) => api.delete(`/users/${id}`))
    );
    toast.success(`${selectedIds.value.length} user berhasil dinonaktifkan`);
    selectedIds.value = [];
    showBulkDeleteModal.value = false;
    fetchData();
  } catch {
    toast.error("Gagal menonaktifkan user");
  }
};

const resetPassword = async (user: User) => {
  resettingPassword.value = true;
  try {
    const response = await api.post(`/users/${user.id}/reset-password`);
    if (response.data.success) {
      resetPasswordResult.value = {
        username: response.data.username,
        new_password: response.data.new_password,
      };
      showResetPasswordModal.value = true;
      toast.success("Password berhasil direset");
    }
  } catch {
    toast.error("Gagal mereset password");
  } finally {
    resettingPassword.value = false;
  }
};

onMounted(() => fetchData());
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Manajemen User</h1>
      <button
        @click="openCreate"
        class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        + Tambah User
      </button>
    </div>

    <DataTable
      :columns="columns"
      :data="data"
      :loading="loading"
      :current-page="currentPage"
      :total-pages="totalPages"
      :total-items="totalItems"
      :per-page="perPage"
      selectable
      searchable
      search-placeholder="Cari username atau nama..."
      :selected-ids="selectedIds"
      @selection-change="onSelectionChange"
      @page-change="onPageChange"
      @per-page-change="onPerPageChange"
      @search="onSearch"
    >
      <template #filters>
        <Dropdown
          v-model="filterRole"
          :options="filterRoleOptions"
          placeholder="Filter Role"
          class="w-40"
          @update:model-value="onFilterRole"
        />
      </template>
      <template #bulk-actions>
        <button
          @click="showBulkDeleteModal = true"
          class="px-3 py-1 bg-red-600 text-white text-sm rounded hover:bg-red-700"
        >
          🚫 Nonaktifkan Terpilih
        </button>
      </template>
      <template #avatar="{ row }">
        <div
          class="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden"
        >
          <img
            v-if="(row as User).avatar_path"
            :src="getAvatarUrl((row as User).avatar_path)"
            class="w-full h-full object-cover"
          />
          <span v-else class="text-gray-500 text-lg">{{
            (row as User).name?.charAt(0)?.toUpperCase()
          }}</span>
        </div>
      </template>
      <template #role="{ row }">
        <span
          class="px-2 py-1 text-xs rounded-full"
          :class="{
          'bg-purple-100 text-purple-700': (row as User).role === 'super_admin',
          'bg-blue-100 text-blue-700': (row as User).role === 'admin',
          'bg-gray-100 text-gray-700': (row as User).role === 'operator',
        }"
          >{{ roleLabels[(row as User).role] }}</span
        >
      </template>
      <template #unit_kerja="{ row }">{{
        (row as User).unit_kerja?.nama || "-"
      }}</template>
      <template #pptk_list="{ row }">
        <div
          v-if="(row as User).pptk_list?.length"
          class="flex flex-wrap gap-1"
        >
          <span
            v-for="up in (row as User).pptk_list?.slice(0, 2)"
            :key="up.id"
            class="px-1.5 py-0.5 bg-blue-100 text-blue-700 text-xs rounded"
          >
            {{ up.pptk?.nama }}
          </span>
          <span
            v-if="((row as User).pptk_list?.length || 0) > 2"
            class="px-1.5 py-0.5 bg-gray-100 text-gray-600 text-xs rounded"
            :title="(row as User).pptk_list?.map(up => up.pptk?.nama).join(', ')"
          >
            +{{ ((row as User).pptk_list?.length || 0) - 2 }} lagi
          </span>
        </div>
        <span v-else class="text-gray-400">-</span>
      </template>
      <template #status="{ row }">
        <span
          class="px-2 py-1 text-xs rounded-full"
          :class="(row as User).is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
        >
          {{ (row as User).is_active ? "Aktif" : "Nonaktif" }}
        </span>
      </template>
      <template #actions="{ row }">
        <div class="flex gap-2">
          <button
            @click.stop="openAvatarUpload(row as User)"
            class="text-green-600 hover:text-green-800"
            title="Upload Avatar"
          >
            📷
          </button>
          <button
            @click.stop="openEdit(row as User)"
            class="text-blue-600 hover:text-blue-800"
            title="Edit User"
          >
            Edit
          </button>
          <button
            @click.stop="resetPassword(row as User)"
            class="text-orange-600 hover:text-orange-800"
            title="Reset Password"
            :disabled="resettingPassword"
          >
            🔑
          </button>
          <button
            @click.stop="toggleActive(row as User)"
            :class="(row as User).is_active ? 'text-red-600 hover:text-red-800' : 'text-green-600 hover:text-green-800'"
          >
            {{ (row as User).is_active ? "Nonaktifkan" : "Aktifkan" }}
          </button>
        </div>
      </template>
    </DataTable>

    <!-- Form Modal -->
    <Modal
      :show="showModal"
      :title="editingItem ? 'Edit User' : 'Tambah User'"
      size="lg"
      @close="showModal = false"
    >
      <form @submit.prevent="save" class="space-y-4">
        <InputField
          v-model="form.username"
          label="Username"
          :error="errors.username"
          required
          :disabled="!!editingItem"
        />
        <InputField
          v-model="form.name"
          label="Nama Lengkap"
          :error="errors.name"
          required
        />
        <InputField
          v-model="form.password"
          label="Password"
          type="password"
          :error="errors.password"
          :required="!editingItem"
          :placeholder="editingItem ? 'Kosongkan jika tidak diubah' : ''"
        />
        <Dropdown
          v-model="form.role"
          :options="roleOptions"
          label="Role"
          :error="errors.role"
          required
        />
        <template v-if="form.role === 'operator'">
          <Dropdown
            v-model="form.unit_kerja_id"
            :options="unitKerjaOptions"
            label="Unit Kerja"
            :error="errors.unit_kerja_id"
            searchable
          />
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              PPTK yang Ditangani
            </label>
            <div
              class="border rounded-lg p-3 max-h-48 overflow-y-auto space-y-2"
              :class="errors.pptk_ids ? 'border-red-500' : 'border-gray-300'"
            >
              <div
                v-for="pptk in filteredPptkOptions"
                :key="pptk.value"
                class="flex items-center gap-2"
              >
                <input
                  type="checkbox"
                  :id="`pptk-${pptk.value}`"
                  :value="pptk.value"
                  v-model="form.pptk_ids"
                  class="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <label
                  :for="`pptk-${pptk.value}`"
                  class="text-sm text-gray-700 cursor-pointer"
                >
                  {{ pptk.label }}
                </label>
              </div>
              <p
                v-if="filteredPptkOptions.length === 0"
                class="text-sm text-gray-500"
              >
                Tidak ada PPTK tersedia
              </p>
            </div>
            <p v-if="errors.pptk_ids" class="mt-1 text-sm text-red-500">
              {{ errors.pptk_ids }}
            </p>
            <p class="mt-1 text-xs text-gray-500">
              Pilih satu atau lebih PPTK yang akan ditangani operator ini
            </p>
          </div>
        </template>
      </form>
      <template #footer>
        <button
          @click="showModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="save"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Simpan
        </button>
      </template>
    </Modal>

    <!-- Avatar Upload Modal -->
    <Modal
      :show="showAvatarModal"
      title="Upload Avatar"
      @close="showAvatarModal = false"
    >
      <div class="flex flex-col items-center gap-4">
        <div
          class="w-32 h-32 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden border-4 border-gray-300"
        >
          <img
            v-if="avatarPreview"
            :src="avatarPreview"
            class="w-full h-full object-cover"
          />
          <span v-else class="text-gray-400 text-4xl">{{
            avatarUser?.name?.charAt(0)?.toUpperCase()
          }}</span>
        </div>
        <input
          type="file"
          accept="image/*"
          @change="onAvatarChange"
          class="text-sm"
        />
        <p class="text-sm text-gray-500">Format: JPG, PNG. Max 2MB</p>
      </div>
      <template #footer>
        <button
          @click="showAvatarModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="uploadAvatar"
          :disabled="!avatarFile"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          Upload
        </button>
      </template>
    </Modal>

    <!-- Delete Confirmation -->
    <Modal
      :show="showDeleteModal"
      title="Konfirmasi Nonaktifkan"
      @close="showDeleteModal = false"
    >
      <p>
        Yakin ingin menonaktifkan user <strong>{{ deletingItem?.name }}</strong
        >?
      </p>
      <template #footer>
        <button
          @click="showDeleteModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="confirmDelete"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Nonaktifkan
        </button>
      </template>
    </Modal>

    <!-- Bulk Delete Confirmation -->
    <Modal
      :show="showBulkDeleteModal"
      title="Konfirmasi Nonaktifkan Massal"
      @close="showBulkDeleteModal = false"
    >
      <p>
        Yakin ingin menonaktifkan <strong>{{ selectedIds.length }}</strong> user
        yang dipilih?
      </p>
      <template #footer>
        <button
          @click="showBulkDeleteModal = false"
          class="px-4 py-2 border rounded-lg hover:bg-gray-50"
        >
          Batal
        </button>
        <button
          @click="confirmBulkDeactivate"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Nonaktifkan Semua
        </button>
      </template>
    </Modal>

    <!-- Reset Password Result Modal -->
    <Modal
      :show="showResetPasswordModal"
      title="Password Berhasil Direset"
      @close="showResetPasswordModal = false"
    >
      <div class="text-center space-y-4">
        <div class="text-6xl">🔑</div>
        <p class="text-gray-600">Password baru untuk user <strong>{{ resetPasswordResult?.username }}</strong>:</p>
        <div class="bg-gray-100 p-4 rounded-lg">
          <p class="text-2xl font-mono font-bold text-blue-600 select-all">{{ resetPasswordResult?.new_password }}</p>
        </div>
        <p class="text-sm text-red-500">⚠️ Catat password ini! Tidak akan ditampilkan lagi.</p>
      </div>
      <template #footer>
        <button
          @click="showResetPasswordModal = false"
          class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Tutup
        </button>
      </template>
    </Modal>
  </div>
</template>
