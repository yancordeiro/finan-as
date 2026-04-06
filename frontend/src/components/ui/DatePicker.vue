<script setup>
import { ref, computed, watch } from 'vue'
import { PopoverRoot, PopoverTrigger, PopoverPortal, PopoverContent } from 'radix-vue'
import { Calendar, ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  modelValue: String, // Formato: YYYY-MM-DD
  label: String,
  error: String,
  disabled: Boolean,
  required: Boolean,
  placeholder: {
    type: String,
    default: 'Selecione uma data'
  }
})

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const selectedDate = ref(props.modelValue || '')

// Calendário
const today = new Date()
const currentMonth = ref(selectedDate.value ? new Date(selectedDate.value + 'T00:00:00') : new Date())

// Nomes dos meses
const monthNames = [
  'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
  'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
]

const weekDays = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']

// Formatar data para exibição
const formatDisplayDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString + 'T00:00:00')
  return date.toLocaleDateString('pt-BR')
}

const displayValue = computed(() => formatDisplayDate(selectedDate.value))

// Gerar dias do mês
const calendarDays = computed(() => {
  const year = currentMonth.value.getFullYear()
  const month = currentMonth.value.getMonth()

  // Primeiro dia do mês
  const firstDay = new Date(year, month, 1)
  // Último dia do mês
  const lastDay = new Date(year, month + 1, 0)

  // Dia da semana do primeiro dia (0 = domingo)
  const startingDayOfWeek = firstDay.getDay()

  const days = []

  // Dias do mês anterior (para preencher a primeira semana)
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = startingDayOfWeek - 1; i >= 0; i--) {
    days.push({
      day: prevMonthLastDay - i,
      date: new Date(year, month - 1, prevMonthLastDay - i),
      isCurrentMonth: false,
      isToday: false,
      isSelected: false
    })
  }

  // Dias do mês atual
  for (let day = 1; day <= lastDay.getDate(); day++) {
    const date = new Date(year, month, day)
    const dateString = formatDateToISO(date)
    const isToday = isSameDay(date, today)
    const isSelected = selectedDate.value === dateString

    days.push({
      day,
      date,
      dateString,
      isCurrentMonth: true,
      isToday,
      isSelected
    })
  }

  // Dias do próximo mês (para preencher a última semana)
  const remainingDays = 42 - days.length // 6 semanas * 7 dias
  for (let day = 1; day <= remainingDays; day++) {
    days.push({
      day,
      date: new Date(year, month + 1, day),
      isCurrentMonth: false,
      isToday: false,
      isSelected: false
    })
  }

  return days
})

// Helpers
const formatDateToISO = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const isSameDay = (date1, date2) => {
  return date1.getDate() === date2.getDate() &&
         date1.getMonth() === date2.getMonth() &&
         date1.getFullYear() === date2.getFullYear()
}

// Navegação
const previousMonth = () => {
  currentMonth.value = new Date(
    currentMonth.value.getFullYear(),
    currentMonth.value.getMonth() - 1,
    1
  )
}

const nextMonth = () => {
  currentMonth.value = new Date(
    currentMonth.value.getFullYear(),
    currentMonth.value.getMonth() + 1,
    1
  )
}

const selectDate = (day) => {
  if (!day.isCurrentMonth) {
    // Se clicou em dia de outro mês, muda o mês e seleciona
    currentMonth.value = new Date(day.date)
  }
  selectedDate.value = day.dateString || formatDateToISO(day.date)
  emit('update:modelValue', selectedDate.value)
  open.value = false
}

const selectToday = () => {
  const todayString = formatDateToISO(today)
  selectedDate.value = todayString
  currentMonth.value = new Date(today)
  emit('update:modelValue', todayString)
  open.value = false
}

// Watch para atualizar quando o modelValue mudar externamente
watch(() => props.modelValue, (newValue) => {
  selectedDate.value = newValue || ''
  if (newValue) {
    currentMonth.value = new Date(newValue + 'T00:00:00')
  }
})
</script>

<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-danger ml-1">*</span>
    </label>

    <PopoverRoot v-model:open="open">
      <PopoverTrigger as-child>
        <button
          type="button"
          :disabled="disabled"
          :class="[
            'input w-full flex items-center justify-between gap-2 text-left',
            error && 'border-danger focus:ring-danger',
            disabled && 'opacity-50 cursor-not-allowed',
            !disabled && 'cursor-pointer hover:border-border-hover'
          ]"
        >
          <span :class="displayValue ? 'text-text-primary' : 'text-text-muted'">
            {{ displayValue || placeholder }}
          </span>
          <Calendar :size="18" class="text-text-secondary flex-shrink-0" />
        </button>
      </PopoverTrigger>

      <PopoverPortal>
        <PopoverContent
          :side="'bottom'"
          :align="'start'"
          :side-offset="8"
          class="z-50 w-auto rounded-lg border bg-card shadow-lg outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95"
        >
          <div class="p-4 space-y-4">
            <!-- Header com navegação de mês -->
            <div class="flex items-center justify-between">
              <button
                type="button"
                @click="previousMonth"
                class="inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground h-8 w-8 transition-colors"
              >
                <ChevronLeft :size="16" />
              </button>

              <div class="text-sm font-semibold">
                {{ monthNames[currentMonth.getMonth()] }} {{ currentMonth.getFullYear() }}
              </div>

              <button
                type="button"
                @click="nextMonth"
                class="inline-flex items-center justify-center rounded-md text-sm font-medium hover:bg-accent hover:text-accent-foreground h-8 w-8 transition-colors"
              >
                <ChevronRight :size="16" />
              </button>
            </div>

            <!-- Dias da semana -->
            <div class="grid grid-cols-7 gap-1">
              <div
                v-for="day in weekDays"
                :key="day"
                class="text-center text-xs font-medium text-text-muted py-2"
              >
                {{ day }}
              </div>
            </div>

            <!-- Dias do mês -->
            <div class="grid grid-cols-7 gap-1">
              <button
                v-for="(day, index) in calendarDays"
                :key="index"
                type="button"
                @click="selectDate(day)"
                :class="[
                  'h-9 w-9 text-sm rounded-md transition-colors',
                  day.isCurrentMonth
                    ? 'text-text-primary hover:bg-accent'
                    : 'text-text-muted hover:bg-accent/50',
                  day.isToday && 'font-bold border border-primary',
                  day.isSelected && 'bg-primary text-white hover:bg-primary/90',
                  !day.isSelected && !day.isToday && 'hover:bg-accent'
                ]"
              >
                {{ day.day }}
              </button>
            </div>

            <!-- Footer com botão "Hoje" -->
            <div class="pt-2 border-t">
              <button
                type="button"
                @click="selectToday"
                class="w-full btn btn-outline text-xs py-1.5"
              >
                Hoje
              </button>
            </div>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>

    <p v-if="error" class="mt-1 text-sm text-danger">
      {{ error }}
    </p>
  </div>
</template>
