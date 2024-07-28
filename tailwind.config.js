/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['internal/view/*.templ'],
  theme: {},
  plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography')]
};
