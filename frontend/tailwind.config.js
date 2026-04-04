/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        bg: {
          DEFAULT: '#1a1a2e',
          light: '#1e1e35',
        },
        card: {
          DEFAULT: '#2d2d44',
          hover: '#35354f',
          border: '#3a3a55',
        },
        primary: {
          DEFAULT: '#8B5CF6',
          hover: '#7c3aed',
          light: '#a78bfa',
          muted: '#4c3d8f',
        },
        text: {
          DEFAULT: '#e8e8e8',
          muted: '#9ca3af',
          secondary: '#6b7280',
          dim: '#4b4b6a',
        },
        success: '#10b981',
        danger: '#ef4444',
        warning: '#f59e0b',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.4s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(16px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
