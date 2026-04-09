<template>
  <div
    class="relative overflow-hidden rounded-2xl min-h-[calc(100vh-8rem)] flex items-center justify-center"
    style="background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)"
  >
    <!-- Dot grid background -->
    <div class="absolute inset-0 opacity-[0.12]" style="background-image: radial-gradient(circle, #94a3b8 1px, transparent 1px); background-size: 28px 28px;"></div>

    <!-- Faint 404 watermark -->
    <span class="absolute inset-0 flex items-center justify-center text-[28vw] font-black text-slate-700/20 select-none pointer-events-none leading-none">404</span>

    <!-- Content -->
    <div class="relative z-10 text-center px-6 max-w-sm">
      <!-- Icon ring -->
      <div class="w-20 h-20 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center mx-auto mb-6 shadow-xl">
        <svg class="w-9 h-9 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z" />
        </svg>
      </div>

      <h1 class="text-3xl font-extrabold text-white mb-2 tracking-tight">Página no encontrada</h1>
      <p class="text-slate-400 text-sm mb-8 leading-relaxed">
        La dirección que buscas no existe o ha sido movida.<br>Verifica la URL o regresa al inicio.
      </p>

      <RouterLink
        to="/"
        class="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl text-sm font-semibold text-white transition-all hover:-translate-y-0.5 hover:shadow-lg"
        style="background: linear-gradient(135deg, #6366f1, #4338ca); box-shadow: 0 4px 14px rgba(99,102,241,0.4)"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Volver al inicio
      </RouterLink>
    </div>
  </div>
</template>

<script setup>
// 404 — Not Found
</script>
