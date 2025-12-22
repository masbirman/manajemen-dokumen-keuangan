<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { InputField } from '@/components/ui'

const authStore = useAuthStore()
const toast = useToast()

const saving = ref(false)
const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const form = reactive({
  name: '',
  username: '',
  email: '',
  password: ''
})

// Initialize form
onMounted(() => {
  if (authStore.user) {
    form.name = authStore.user.name
    form.username = authStore.user.username
    form.email = (authStore.user as any).email || ''
  }
})

const handleAvatarClick = () => {
  fileInput.value?.click()
}

const handleFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (!target.files?.length) return

  const file = target.files[0]
  if (file.size > 2 * 1024 * 1024) {
    toast.error('Ukuran file maksimal 2MB')
    return
  }

  const formData = new FormData()
  formData.append('avatar', file)

  uploading.value = true
  try {
    const response = await api.post('/auth/profile/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    // Update store
    if (authStore.user) {
      if (typeof authStore.user.avatar_path !== 'undefined') {
        authStore.user.avatar_path = response.data.avatar_path
      } else {
        // Force update if property missing
        (authStore.user as any).avatar_path = response.data.avatar_path
      }
    }
    
    toast.success('Avatar berhasil diperbarui')
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal upload avatar')
  } finally {
    uploading.value = false
    // Reset input
    if (fileInput.value) fileInput.value.value = ''
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    const payload: any = {
      name: form.name,
      username: form.username
    }
    
    if (form.password) {
      payload.password = form.password
    }

    const response = await api.put('/auth/profile', payload)
    
    // Update store
    if (authStore.user) {
        Object.assign(authStore.user, response.data.data)
    }
    
    form.password = '' // Clear password field
    toast.success('Profil berhasil diperbarui')
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Gagal menyimpan profil')
  } finally {
    saving.value = false
  }
}

// Helper for avatar URL
const apiBaseUrl = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8000'
const getAvatarUrl = (path: string | undefined) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return apiBaseUrl + path
}
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <div class="bg-white rounded-xl shadow-sm border border-secondary-200 overflow-hidden">
        <!-- Header -->
        <div class="p-6 border-b border-secondary-100 bg-secondary-50">
            <h1 class="text-xl font-bold text-secondary-900">Edit Profil</h1>
            <p class="text-sm text-secondary-500">Perbarui informasi akun Anda</p>
        </div>

        <div class="p-6">
            <!-- Avatar Section -->
            <div class="flex flex-col items-center mb-8">
                <div class="relative group cursor-pointer" @click="handleAvatarClick">
                    <div class="w-24 h-24 rounded-full overflow-hidden border-4 border-white shadow-md bg-secondary-100">
                        <img 
                            v-if="authStore.user?.avatar_path" 
                            :src="getAvatarUrl(authStore.user.avatar_path)" 
                            alt="Avatar" 
                            class="w-full h-full object-cover"
                        />
                         <div v-else class="w-full h-full flex items-center justify-center text-secondary-400 text-3xl font-bold bg-secondary-200">
                            {{ authStore.user?.name?.charAt(0).toUpperCase() }}
                        </div>
                    </div>
                    
                    <!-- Overlay -->
                    <div class="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <svg v-if="!uploading" class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                        </svg>
                        <svg v-else class="w-8 h-8 text-white animate-spin" fill="none" viewBox="0 0 24 24">
                             <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                             <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                    </div>
                </div>
                <p class="text-sm text-secondary-500 mt-2">Klik gambar untuk mengubah avatar</p>
                <input 
                    type="file" 
                    ref="fileInput" 
                    class="hidden" 
                    accept="image/*"
                    @change="handleFileChange"
                />
            </div>

            <!-- Form -->
            <form @submit.prevent="saveProfile" class="space-y-4">
                <InputField 
                    v-model="form.name" 
                    label="Nama Lengkap" 
                    placeholder="Masukkan nama lengkap"
                    required
                />
                
                <InputField 
                    v-model="form.username" 
                    label="Username" 
                    placeholder="Masukkan username"
                    required
                />

                <InputField 
                    v-model="form.email" 
                    label="Email" 
                    type="email"
                    placeholder="Masukkan email (opsional)"
                />

                <div class="border-t border-secondary-100 pt-4 mt-6">
                    <h3 class="text-sm font-semibold text-secondary-900 mb-4">Ganti Password</h3>
                    <InputField 
                        v-model="form.password" 
                        label="Password Baru" 
                        type="password"
                        placeholder="Biarkan kosong jika tidak ingin mengubah"
                    />
                     <p class="text-xs text-secondary-500 mt-1">Minimal 6 karakter.</p>
                </div>

                <div class="flex justify-end pt-6">
                    <button 
                        type="submit" 
                        :disabled="saving"
                        class="btn-primary flex items-center gap-2"
                    >
                        <svg v-if="saving" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <span>{{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
                    </button>
                </div>
            </form>
        </div>
    </div>
  </div>
</template>
