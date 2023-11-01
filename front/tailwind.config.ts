import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic':
          'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
      },
      colors: {
        primary: {
          default: '#f9842c',
          dark: '#FA6C28',
          light: '#f9842c'
        },
      }
    },
    fontFamily: {
      Nunito: ["var(--font-Nunito)"],
      Noto: ["var(--font-Noto)"],
    }
  },
  plugins: [],
}
export default config
