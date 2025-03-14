import React from "react";
import { View, Text, Image, TouchableOpacity, ScrollView } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";

const TEMP_EVENTS = [
  {
    id: "1",
    title: "Благодійний марафон",
    description: `Долучайтеся до нашого благодійного марафону, де ви зможете взяти участь у захопливих активностях, підтримати добру справу та отримати незабутні враження!`,
    image:
      "https://fartlek.com.ua/wp-content/uploads/2022/08/balakliia-run-2022-1140x640.jpg",
    start_date: "2025-04-01",
    end_date: "2025-04-10",
    format: "offline",
    minimum_donation: 50,
    max_participants: 100,
    location: "Київ, Хрещатик 1",
  },
  {
    id: "2",
    title: "Онлайн-лекція з мотивації",
    description: `Запрошуємо на унікальну лекцію, де ви дізнаєтеся секрети успіху від відомого спікера. Це ваш шанс змінити життя на краще!`,
    image: "https://i.ytimg.com/vi/s-0V76nC1S4/maxresdefault.jpg",
    start_date: "2025-05-15",
    end_date: "2025-05-15",
    format: "online",
    minimum_donation: 0,
    max_participants: 500,
    location: "Zoom",
  },
];

export default function DetailedEventCard() {
  const { id } = useLocalSearchParams();
  const router = useRouter();

  const event = TEMP_EVENTS.find((item) => item.id === id);

  if (!event) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <Text className="text-lg text-gray-msg">Подію не знайдено 😕</Text>
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
          source={{ uri: event.image }}
          className="w-full h-96 rounded-lg mb-4"
          resizeMode="cover"
        />

        <Text className="text-xl font-bold text-black text-center mb-2">
          {event.title}
        </Text>

        <Text className="text-gray-msg mt-2 text-base leading-5">
          {event.description}
        </Text>

        <View className="mt-4 bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-grey-msg">
            📅 {event.start_date} – {event.end_date}
          </Text>
          <Text className="text-grey-msg">
            📍 {event.format === "offline" ? event.location : "Онлайн"}
          </Text>
          <Text className="text-grey-msg">
            👥 Макс. учасників: {event.max_participants}
          </Text>
          <Text className="text-grey-msg">
            💰 Мін. внесок:{" "}
            {event.minimum_donation === 0
              ? "Безкоштовно"
              : `${event.minimum_donation} грн`}
          </Text>
        </View>

        <TouchableOpacity className="bg-primary rounded-lg p-3 mt-5 items-center">
          <Text className="text-white text-lg font-semibold">
            Взяти участь 🎟️
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
