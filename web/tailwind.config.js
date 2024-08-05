/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./static/**/*.html'],
  theme: {
    extend: {
      colors: {
        'fft-bg': '#39363d',
        'fft-r-name': '#e0c6dc',
        'fft-txt': '#fffcf9',
        'fft-li': '#95764c',
        'fft-li-num': '#e0c7a5',
        'fft-source': '#95614c'
      },
    },
    fontFamily: {
      'serif': ['"Gilda Display"'],
      'normal-text': ['Maitree'],
    },
  },
  plugins: [],
}

