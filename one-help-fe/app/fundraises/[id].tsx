import React from "react";
import { View, Text, Image, TouchableOpacity, ScrollView } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import { ProgressBar } from "react-native-paper";

const TEMP_FUNDRAISERS = [
  {
    id: "1",
    title: "Збір на 2 ecoflow delta max",
    description: `Наші прикордонники потребують допомоги в забезпеченні та підготовці до бойового виїзду. Їм необхідно:
    • Комплект гуми – 25 500 грн
    • 2 автономні дизельні обігрівачі – 12 700 грн
    • EcoFlow Delta Max 1600 – 38 000 грн

    Хлопці щодня роблять усе можливе, щоб ми могли спокійно спати в теплих ліжках. Не залишаймося байдужими, коли вони потребують нашої підтримки!`,
    image:
      "https://www.peoplesproject.com/wp-content/uploads/2023/09/11113.jpg",
    goal: 71000,
    raised: 55000,
  },
  {
    id: "2",
    title: "Збір на авто",
    description: "25 бригада, Покровського напрямку",
    image:
      "https://ptv.ua/mediafiles/gallery/d5d26a74-41b8-4985-bccb-492ffa84f7b5.jfif",
    goal: 260000,
    raised: 100000,
  },
];

export default function DetailedFundraiseCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();

  const fundraiser = TEMP_FUNDRAISERS.find((item) => item.id === id);

  if (!fundraiser) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <Text className="text-lg text-gray-msg">Збір не знайдено 😕</Text>
        <TouchableOpacity
          className="mt-4 px-4 py-2 bg-primary rounded-md"
          onPress={() => router.back()}
        >
          <Text className="text-white">Назад</Text>
        </TouchableOpacity>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-4 rounded-2xl">
        <Image
          source={{ uri: fundraiser.image }}
          className="w-full h-96 rounded-lg mb-4"
          resizeMode="cover"
        />

        <Text className="text-xl font-bold text-black text-center mb-2">
          {fundraiser.title}
        </Text>

        <Text className="text-gray-msg mt-2 text-base leading-5 mb-50">
          {fundraiser.description}
        </Text>

        <View className="mt-4">
          <Text className="text-grey-msg text-sm text-center font-medium">
            {fundraiser.raised.toLocaleString()} /{" "}
            {fundraiser.goal.toLocaleString()}
          </Text>
          <ProgressBar
            progress={fundraiser.raised / fundraiser.goal}
            color="#2563EB"
            className="h-3 rounded-md mt-1"
          />
        </View>

        <TouchableOpacity className="bg-primary rounded-lg p-3 mt-5 items-center">
          <Text className="text-white text-lg font-semibold">Донат 🍩</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
