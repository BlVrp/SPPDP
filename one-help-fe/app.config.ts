import 'dotenv/config';

export default {
  expo: {
    name: "OneHelp",
    slug: "onehelp",
    extra: {
      API_URL: process.env.API_URL,
    },
  },
};
