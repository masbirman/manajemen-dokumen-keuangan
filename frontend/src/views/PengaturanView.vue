<script setup lang="ts">
import { ref, onMounted } from "vue";
import api from "@/services/api";
import { useToast } from "@/composables/useToast";
import { InputField } from "@/components/ui";

interface Settings {
  app_name: string;
  app_description: string;
  organization_name: string;
  organization_address: string;
}

const toast = useToast();
const loading = ref(false);
const saving = ref(false);
const settings = ref<Settings>({
  app_name: "",
  app_description: "",
  organization_name: "",
  organization_address: "",
});

const fetchSettings = async () => {
  loading.value = true;
  try {
    const response = await api.get("/settings");
    const data = response.data.data || response.data || {};
    settings.value = {
      app_name: data.app_name || "",
      app_description: data.app_description || "",
      organization_name: data.organization_name || "",
      organization_address: data.organization_address || "",
    };
  } catch {
    toast.error("Gagal memuat pengaturan");
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  saving.value = true;
  try {
    await api.put("/settings", settings.value);
    toast.success("Pengaturan berhasil disimpan");
  } catch {
    toast.error("Gagal menyimpan pengaturan");
  } finally {
    saving.value = false;
  }
};

onMounted(fetchSettings);
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Pengaturan</h1>

    <div
      v-if="loading"
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center"
    >
      <p class="text-gray-500">Loading...</p>
    </div>

    <div
      v-else
      class="bg-white rounded-xl shadow-sm p-6 border border-gray-100"
    >
      <form @submit.prevent="saveSettings" class="space-y-4">
        <h2 class="text-lg font-semibold text-gray-700 mb-4">
          Informasi Aplikasi
        </h2>
        <InputField v-model="settings.app_name" label="Nama Aplikasi" />
        <InputField
          v-model="settings.app_description"
          label="Deskripsi Aplikasi"
        />

        <h2 class="text-lg font-semibold text-gray-700 mb-4 mt-6">
          Informasi Organisasi
        </h2>
        <InputField
          v-model="settings.organization_name"
          label="Nama Organisasi"
        />
        <InputField
          v-model="settings.organization_address"
          label="Alamat Organisasi"
        />

        <div class="flex justify-end pt-4">
          <button
            type="submit"
            :disabled="saving"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ saving ? "Menyimpan..." : "Simpan Pengaturan" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
