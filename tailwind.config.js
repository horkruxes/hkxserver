module.exports = {
  purge: {
    enabled: ((process.env.ENV === 'production') ? true : false),
    content: ["./templates/**/*.html"],
  },
  darkMode: 'media', // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
