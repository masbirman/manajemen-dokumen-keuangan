<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { useAuthStore } from "@/stores/auth";
import {
  InputField,
  Dropdown,
  CurrencyInput,
  FileUpload,
} from "@/components/ui";

interface Option {
  id: string;
  nama: string;
  kode?: string;
  unit_kerja_id?: string;
}

const toast = useToast();
const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const loadingData = ref(true);
const unitKerjaList = ref<Option[]>([]);
const pptkList = ref<Option[]>([]);
const sumberDanaList = ref<Option[]>([]);
const jenisDokumenList = ref<Option[]>([]);

const form = ref({
  nomor_dokumen: "",
  tanggal_dokumen: "",
  uraian: "",
  nilai: 0,
  nomor_kwitansi: "",
  unit_kerja_id: "",
  pptk_id: "",
  sumber_dana_id: "",
  jenis_dokumen_id: "",
});
const file = ref<File | null>(null);
const currentFilePath = ref("");
const errors = ref<Record<string, string>>({});

const dokumenId = computed(() => route.params.id as string);

const unitKerjaOptions = computed(() =>
  unitKerjaList.value.map((u) => ({ value: u.id, label: u.nama }))
);
const pptkOptions = computed(() => {
  let filtered = form.value.unit_kerja_id
    ? pptkList.value.filter((p) => p.unit_kerja_id === form.value.unit_kerja_id)
    : pptkList.value;

  // For operators with assigned PPTKs, filter to only show their assigned PPTKs
  if (authStore.isOperator && authStore.user?.pptk_list?.length) {
    const assignedPPTKIds = authStore.user.pptk_list.map((up) => up.pptk_id);
    filtered = filtered.filter((p) => assignedPPTKIds.includes(p.id));
  }

  return filtered.map((p) => ({ value: p.id, label: p.nama }));
});
const sumberDanaOptions = computed(() =>
  sumberDanaList.value.map((s) => ({
    value: s.id,
    label: `${s.kode} - ${s.nama}`,
  }))
);
const jenisDokumenOptions = computed(() =>
  jenisDokumenList.value.map((j) => ({
    value: j.id,
    label: `${j.kode} - ${j.nama}`,
  }))
);

const fetchDokumen = async () => {
  try {
    const response = await api.get(`/dokumen/${dokumenId.value}`);
    const dokumen = response.data.data;

    form.value = {
      nomor_dokumen: dokumen.nomor_dokumen || "",
      tanggal_dokumen: dokumen.tanggal_dokumen
        ? dokumen.tanggal_dokumen.split("T")[0]
        : "",
      uraian: dokumen.uraian || "",
      nilai: dokumen.nilai || 0,
      nomor_kwitansi: dokumen.nomor_kwitansi || "",
      unit_kerja_id: dokumen.unit_kerja_id || "",
      pptk_id: dokumen.pptk_id || "",
      sumber_dana_id: dokumen.sumber_dana_id || "",
      jenis_dokumen_id: dokumen.jenis_dokumen_id || "",
    };
    currentFilePath.value = dokumen.file_path || "";
  } catch {
    toast.error("Gagal memuat data dokumen");
    router.push("/dokumen");
  }
};

const fetchMasterData = async () => {
  try {
    const [ukRes, pptkRes, sdRes, jdRes] = await Promise.all([
      api.get("/unit-kerja/active"),
      api.get("/pptk/active"),
      api.get("/sumber-dana/active"),
      api.get("/jenis-dokumen/active"),
    ]);
    unitKerjaList.value = ukRes.data.data || ukRes.data || [];
    pptkList.value = pptkRes.data.data || pptkRes.data || [];
    sumberDanaList.value = sdRes.data.data || sdRes.data || [];
    jenisDokumenList.value = jdRes.data.data || jdRes.data || [];
  } catch {
    toast.error("Gagal memuat master data");
  }
};

const MAX_FILE_SIZE = 300 * 1024; // 300KB

const handleFileUpdate = (f: File | null) => {
  if (f && f.size > MAX_FILE_SIZE) {
    errors.value.file = `Ukuran file terlalu besar (${(f.size / 1024).toFixed(
      0
    )}KB). Maksimal 300KB`;
    file.value = null;
    toast.error("Ukuran file melebihi batas maksimal 300KB");
    return;
  }
  errors.value.file = "";
  file.value = f;
};

const submit = async () => {
  errors.value = {};
  loading.value = true;

  const formData = new FormData();
  Object.entries(form.value).forEach(([key, value]) => {
    formData.append(key, String(value));
  });
  if (file.value) {
    formData.append("file", file.value);
  }

  try {
    await api.put(`/dokumen/${dokumenId.value}`, formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    toast.success("Dokumen berhasil diupdate");
    router.push("/dokumen");
  } catch (e: unknown) {
    const err = e as {
      response?: { data?: { errors?: Record<string, string>; error?: string } };
    };
    errors.value = err.response?.data?.errors || {};
    toast.error(err.response?.data?.error || "Gagal mengupdate dokumen");
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  loadingData.value = true;
  await fetchMasterData();
  await fetchDokumen();
  loadingData.value = false;
});
</script>

<template>
  <div class="max-w-3xl">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Edit Dokumen</h1>

    <div v-if="loadingData" class="text-center py-8">
      <p class="text-gray-500">Memuat data...</p>
    </div>

    <div
      v-else
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100"
    >
      <form @submit.prevent="submit" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <InputField
            v-model="form.nomor_dokumen"
            label="Nomor Dokumen"
            :error="errors.nomor_dokumen"
            readonly
            required
          />
          <InputField
            v-model="form.tanggal_dokumen"
            label="Tanggal Dokumen"
            type="date"
            :error="errors.tanggal_dokumen"
            readonly
            required
          />
        </div>

        <InputField
          v-model="form.uraian"
          label="Uraian"
          :error="errors.uraian"
          required
        />

        <div class="grid grid-cols-2 gap-4">
          <InputField
            v-model="form.nomor_kwitansi"
            label="Nomor Kwitansi / Nota Pesanan"
            :error="errors.nomor_kwitansi"
            placeholder="Masukkan nomor kwitansi atau nota pesanan"
          />
          <CurrencyInput
            v-model="form.nilai"
            label="Nilai (Rp)"
            :error="errors.nilai"
            required
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <Dropdown
            v-model="form.unit_kerja_id"
            :options="unitKerjaOptions"
            label="Unit Kerja"
            :error="errors.unit_kerja_id"
            :disabled="authStore.isOperator && !!authStore.user?.unit_kerja_id"
            searchable
            required
          />
          <Dropdown
            v-model="form.pptk_id"
            :options="pptkOptions"
            label="PPTK"
            :error="errors.pptk_id"
            :disabled="
              authStore.isOperator &&
              !!authStore.user?.pptk_id &&
              !authStore.user?.pptk_list?.length
            "
            searchable
            required
          />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <Dropdown
            v-model="form.sumber_dana_id"
            :options="sumberDanaOptions"
            label="Sumber Dana"
            :error="errors.sumber_dana_id"
            searchable
            required
          />
          <Dropdown
            v-model="form.jenis_dokumen_id"
            :options="jenisDokumenOptions"
            label="Jenis Dokumen"
            :error="errors.jenis_dokumen_id"
            searchable
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            File Dokumen (PDF) - Kosongkan jika tidak ingin mengubah
          </label>
          <p class="text-xs text-gray-500 mb-2">Maksimal ukuran file: 300KB</p>
          <p v-if="currentFilePath" class="text-sm text-gray-500 mb-2">
            File saat ini: {{ currentFilePath.split("/").pop() }}
          </p>
          <FileUpload
            accept=".pdf"
            :error="errors.file"
            @update:file="handleFileUpdate"
          />
          <p v-if="file" class="mt-1 text-sm text-green-600">
            ✓ File baru: {{ file.name }} ({{ (file.size / 1024).toFixed(0) }}KB)
          </p>
        </div>

        <div class="flex justify-end gap-3 pt-4">
          <button
            type="button"
            @click="router.back()"
            class="px-6 py-2 border rounded-lg hover:bg-gray-50"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ loading ? "Menyimpan..." : "Simpan" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
