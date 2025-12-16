<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";
import {
  InputField,
  Dropdown,
  CurrencyInput,
  FileUpload,
} from "@/components/ui";
import DocumentScanner from "@/components/DocumentScanner.vue";
import PetunjukDisplay from "@/components/PetunjukDisplay.vue";

interface Option {
  id: string;
  nama: string;
  kode?: string;
  unit_kerja_id?: string;
}

interface Petunjuk {
  id: string;
  judul: string;
  konten: string;
}

const toast = useToast();
const authStore = useAuthStore();
const router = useRouter();

const loading = ref(false);
const showScanner = ref(false);
const unitKerjaList = ref<Option[]>([]);
const pptkList = ref<Option[]>([]);
const sumberDanaList = ref<Option[]>([]);
const jenisDokumenList = ref<Option[]>([]);
const petunjukList = ref<Petunjuk[]>([]);

const form = ref({
  nomor_dokumen: "",
  tanggal_dokumen: new Date().toISOString().split("T")[0],
  uraian: "",
  nilai: 0,
  unit_kerja_id: "",
  pptk_id: "",
  sumber_dana_id: "",
  jenis_dokumen_id: "",
  nomor_kwitansi: "",
});
const file = ref<File | null>(null);
const errors = ref<Record<string, string>>({});

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

// Reset form to initial values
const resetForm = () => {
  form.value = {
    nomor_dokumen: "",
    tanggal_dokumen: new Date().toISOString().split("T")[0],
    uraian: "",
    nilai: 0,
    unit_kerja_id: "",
    pptk_id: "",
    sumber_dana_id: "",
    jenis_dokumen_id: "",
    nomor_kwitansi: "",
  };
  file.value = null;
  errors.value = {};
};

const generateNomorDokumen = (unitKerjaNama: string) => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  return `DOK/${year}${month}${day}/${unitKerjaNama}`;
};

const fetchMasterData = async () => {
  try {
    // Use /active endpoints for dropdowns (returns all active items without pagination)
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

    // Set defaults for operator based on assigned unit_kerja and pptk
    if (authStore.isOperator && authStore.user) {
      if (authStore.user.unit_kerja_id) {
        form.value.unit_kerja_id = authStore.user.unit_kerja_id;
      }
      // If operator has multiple PPTKs, set first one as default
      if (authStore.user.pptk_list?.length) {
        form.value.pptk_id = authStore.user.pptk_list[0].pptk_id;
      } else if (authStore.user.pptk_id) {
        form.value.pptk_id = authStore.user.pptk_id;
      }
    }

    // Auto-generate nomor dokumen with unit kerja nama
    const selectedUnitKerja = unitKerjaList.value.find(
      (u) => u.id === form.value.unit_kerja_id
    );
    const unitKerjaNama = selectedUnitKerja?.nama || "UNIT";
    form.value.nomor_dokumen = generateNomorDokumen(unitKerjaNama);
  } catch (e) {
    console.error("Error fetching master data:", e);
    toast.error("Gagal memuat master data");
  }
};

const fetchPetunjuk = async () => {
  try {
    const response = await api.get("/petunjuk/halaman/input_dokumen");
    petunjukList.value = response.data.data || [];
  } catch {
    // Ignore error for petunjuk
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

const handleScannedPdf = (f: File) => {
  file.value = f;
  showScanner.value = false;
  toast.success("PDF dari scan berhasil dibuat");
};



const submit = async () => {
  errors.value = {};
  if (!file.value) {
    errors.value.file = "File dokumen wajib diupload";
    return;
  }

  loading.value = true;
  const formData = new FormData();
  Object.entries(form.value).forEach(([key, value]) => {
    formData.append(key, String(value));
  });
  formData.append("file", file.value);

  try {
    await api.post("/dokumen", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    toast.success("Dokumen berhasil disimpan");
    router.push("/dokumen");
  } catch (e: unknown) {
    const err = e as {
      response?: { data?: { errors?: Record<string, string>; error?: string } };
    };
    errors.value = err.response?.data?.errors || {};
    toast.error(err.response?.data?.error || "Gagal menyimpan dokumen");
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  resetForm();
  fetchMasterData();
  fetchPetunjuk();
});
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Input Dokumen</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Form Section -->
      <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <form @submit.prevent="submit" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <InputField
              v-model="form.nomor_dokumen"
              label="Nomor Dokumen"
              :error="errors.nomor_dokumen"
              readonly
              required
            />
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Waktu Penginputan
              </label>
              <div class="px-3 py-2 bg-gray-100 border border-gray-200 rounded-lg text-gray-700">
                {{ new Date().toLocaleDateString("id-ID") }} {{ new Date().toLocaleTimeString("id-ID", { hour: '2-digit', minute: '2-digit' }) }}
              </div>
              <p class="text-xs text-gray-400 mt-1">Waktu akan tercatat otomatis saat menyimpan</p>
            </div>
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
              :disabled="
                authStore.isOperator && !!authStore.user?.unit_kerja_id
              "
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
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >File Dokumen (PDF)</label
            >
            
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- Upload Option -->
              <div>
                <FileUpload
                  accept=".pdf"
                  :error="errors.file"
                  @update:file="handleFileUpdate"
                  required
                />
              </div>

              <!-- Scan Option -->
              <button
                type="button"
                @click="showScanner = true"
                class="border-2 border-dashed border-blue-300 bg-blue-50 rounded-xl p-4 flex flex-col items-center justify-center cursor-pointer hover:bg-blue-100 transition-all gap-2 group min-h-[160px]"
              >
                <div class="w-12 h-12 bg-white rounded-full shadow-sm flex items-center justify-center text-2xl group-hover:scale-110 transition-transform">
                  📷
                </div>
                <div class="text-center">
                  <p class="font-semibold text-blue-700">Scan Dokumen</p>
                  <p class="text-xs text-blue-500 mt-1">Gunakan kamera HP/Tablet</p>
                </div>
              </button>
            </div>

            <div class="flex justify-between items-start mt-2">
              <p class="text-xs text-gray-500">
                Maksimal ukuran file: 300KB
              </p>
              <p v-if="file" class="text-sm font-medium text-green-600 flex items-center gap-1">
                <span>✓</span>
                {{ file.name }}
                <span class="text-gray-400 font-normal">({{ (file.size / 1024).toFixed(0) }}KB)</span>
              </p>
            </div>
            <p v-if="errors.file" class="text-xs text-red-500 mt-1">{{ errors.file }}</p>
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

      <!-- Petunjuk Section -->
      <div v-if="petunjukList.length > 0">
        <div
          class="bg-blue-50 rounded-xl shadow-sm p-6 border border-blue-100 sticky top-4"
        >
          <h2
            class="text-lg font-semibold text-blue-800 mb-4 flex items-center gap-2"
          >
            📖 Petunjuk
          </h2>
          <div class="space-y-4">
            <PetunjukDisplay
              v-for="petunjuk in petunjukList"
              :key="petunjuk.id"
              :judul="petunjuk.judul"
              :konten="petunjuk.konten"
            />
          </div>
        </div>
      </div>
    </div>
  </div>

  <DocumentScanner
    v-if="showScanner"
    @pdf-created="handleScannedPdf"
    @close="showScanner = false"
  />
</template>

<style scoped>
.petunjuk-content :deep(h1) {
  font-size: 1.25rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}
.petunjuk-content :deep(h2) {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}
.petunjuk-content :deep(h3) {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}
.petunjuk-content :deep(p) {
  margin-bottom: 0.5rem;
}
.petunjuk-content :deep(ul),
.petunjuk-content :deep(ol) {
  padding-left: 1.5rem;
  margin-bottom: 0.5rem;
}
.petunjuk-content :deep(ul) {
  list-style-type: disc !important;
}
.petunjuk-content :deep(ol) {
  list-style-type: decimal !important;
}
.petunjuk-content :deep(li) {
  margin-bottom: 0.25rem;
  display: list-item !important;
}
.petunjuk-content :deep(a) {
  color: #2563eb;
  text-decoration: underline;
}
.petunjuk-content :deep(strong) {
  font-weight: 600;
}
.petunjuk-content :deep(em) {
  font-style: italic;
}
.petunjuk-content :deep(u) {
  text-decoration: underline;
}
.petunjuk-content :deep(s) {
  text-decoration: line-through;
}
.petunjuk-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin: 0.5rem 0;
}
/* Alignment support */
.petunjuk-content :deep(.ql-align-center),
.petunjuk-content :deep([style*="text-align: center"]) {
  text-align: center !important;
}
.petunjuk-content :deep(.ql-align-right),
.petunjuk-content :deep([style*="text-align: right"]) {
  text-align: right !important;
}
.petunjuk-content :deep(.ql-align-justify),
.petunjuk-content :deep([style*="text-align: justify"]) {
  text-align: justify !important;
}
.petunjuk-content :deep(p[style*="text-align: center"]) img {
  display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>
