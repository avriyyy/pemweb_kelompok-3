tailwind.config = {
  theme: {
    extend: {
      colors: {
        // Brand
        primary: '#c5e384',
        'primary-hover': '#b3d572',

        // Dark surfaces
        dark: '#200f07',
        'dark-soft': '#2b1810',
        'dark-elevated': '#36200f',

        // Light surfaces
        surface: '#ffffff',
        'surface-muted': '#f7f9ec',

        // Borders
        'border-light': '#e3e8d0',

        // Status
        success: '#4ade80',
        danger: '#f87171',
        warning: '#fbbf24',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
};
