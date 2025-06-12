/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.js",
    "./utils/**/*.js",
    "./components/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        'hospedate-green': '#227f4b',
        'hosp-pale-blue': '#59a5bf',
        'dis-hosp-pale-blue': '#bddae4',
        'hosp-light-blue': '#2A8BD2'
      }
    },
    screens: {
      "sm": '744px',
      'md': '744px',
      'desktop': '1440px',
      'laptop': '1128px',
      'tablet': '744px',
    }
  },
  plugins: [
  ],
}
