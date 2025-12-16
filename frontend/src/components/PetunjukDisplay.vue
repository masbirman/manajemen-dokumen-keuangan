<script setup lang="ts">
import { computed } from "vue";

interface InfoItem {
  text: string;
  isBold: boolean;
  linkText?: string;
  linkUrl?: string;
}

interface InfoSection {
  title: string;
  icon: string;
  color: string;
  items: InfoItem[];
}

interface PetunjukContent {
  sections: InfoSection[];
  image_url?: string;
  image_size?: number;
}

const props = defineProps<{
  judul: string;
  konten: string;
}>();

// Parse JSON content
const parsedContent = computed<PetunjukContent | null>(() => {
  if (!props.konten) return null;
  try {
    return JSON.parse(props.konten);
  } catch {
    return null;
  }
});

// Get full image URL
const getImageUrl = (url: string | undefined) => {
  if (!url) return "";
  if (url.startsWith("/")) {
    return `http://localhost:8000${url}`;
  }
  return url;
};

// Format item text with bold and links
const formatItemText = (item: InfoItem): string => {
  let text = item.text;
  
  if (item.isBold) {
    text = `<b>${text}</b>`;
  }
  
  if (item.linkText && text.includes(item.linkText)) {
    if (item.linkUrl) {
      text = text.replace(
        item.linkText,
        `<a href="${item.linkUrl}" target="_blank" class="text-blue-600 hover:underline">${item.linkText}</a>`
      );
    } else {
      text = text.replace(
        item.linkText,
        `<b class="text-blue-600">${item.linkText}</b>`
      );
    }
  }
  
  return text;
};
</script>

<template>
  <div class="bg-white rounded-lg p-4 border border-blue-100">
    
    <!-- Parsed content (new JSON format) -->
    <template v-if="parsedContent">
      <!-- Image if exists -->
      <div v-if="parsedContent.image_url" class="mb-4">
        <img
          :src="getImageUrl(parsedContent.image_url)"
          alt="Petunjuk Image"
          class="rounded-lg object-contain"
          :style="{ maxHeight: (parsedContent.image_size || 200) + 'px' }"
        />
      </div>
      
      <!-- Sections -->
      <div class="space-y-4">
        <div
          v-for="(section, sectionIndex) in parsedContent.sections"
          :key="sectionIndex"
          class="border-l-4 pl-3"
          :style="{ borderColor: section.color }"
        >
          <h4 class="font-semibold text-gray-700 mb-2 flex items-center gap-2 text-sm">
            <span :style="{ color: section.color }">{{ section.icon }}</span>
            {{ section.title }}
          </h4>
          <ul class="text-sm text-gray-600 space-y-1.5">
            <li 
              v-for="(item, itemIndex) in section.items" 
              :key="itemIndex"
              v-html="`${itemIndex + 1}. ${formatItemText(item)}`"
              class="leading-relaxed"
            ></li>
          </ul>
        </div>
      </div>
    </template>
    
    <!-- Fallback: Raw HTML content (old format) -->
    <div
      v-else
      class="text-sm text-gray-600 petunjuk-content"
      v-html="konten"
    ></div>
  </div>
</template>

<style scoped>
.petunjuk-content :deep(b) {
  font-weight: 600;
}
.petunjuk-content :deep(a) {
  color: #2563eb;
  text-decoration: underline;
}
.petunjuk-content :deep(img) {
  max-width: 100%;
  border-radius: 0.5rem;
  margin: 0.5rem 0;
}
</style>
