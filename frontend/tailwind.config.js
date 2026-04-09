/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
        'ice-blue':    '#E6F3FF',
        'nordic-cyan': '#00D4FF',
        'cool-gray':   '#6B7280',
        'dark-gray':   '#374151',
      },
      animation: {
        'fade-in':      'fadeIn .2s ease-out',
        'slide-up':     'slideUp .25s ease-out',
        'overlay-in':   'overlayIn .2s ease-out',
        'spin-slow':    'spin 1.2s linear infinite',
      },
      keyframes: {
        fadeIn: {
          '0%':   { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%':   { opacity: '0', transform: 'translateY(12px) scale(.98)' },
          '100%': { opacity: '1', transform: 'translateY(0) scale(1)' },
        },
        overlayIn: {
          '0%':   { opacity: '0' },
          '100%': { opacity: '1' },
        },
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}