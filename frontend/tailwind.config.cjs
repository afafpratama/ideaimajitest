const withMT = require("@material-tailwind/react/utils/withMT");
 
module.exports = withMT({
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      width: {
        'no-scroll': 'auto',
      },
    },
  },
  variants: {},
  plugins: [
    function ({ addUtilities }) {
      const newUtilities = {
        '.hide-scrollbar': {
          'scrollbar-width': 'none', /* Firefox */
          '-ms-overflow-style': 'none', /* Internet Explorer 10+ */
        },
        '.hide-scrollbar::-webkit-scrollbar': {
          display: 'none', /* Safari and Chrome */
        },
      };
      addUtilities(newUtilities, ['responsive', 'hover']);
    },
  ],
});