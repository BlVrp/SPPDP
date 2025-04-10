import React from "react";
import { View, Text, TouchableOpacity, Image } from "react-native";
import { useRouter } from "expo-router";

export default function EventCard({ event }: { event: any }) {
  const router = useRouter();

  const navigateToEvent = () => {
    router.push({
      pathname: "/events/[id]",
      params: { id: `${event.id}` },
    });
  };

  return (
    <View className="bg-accent p-4 rounded-lg mb-4">
      <TouchableOpacity onPress={navigateToEvent} activeOpacity={0.8}>
        <View className="w-full">
          <Image
            source={{ uri: event.image }}
            className="w-full h-32 rounded-lg mb-3"
            resizeMode="cover"
          />
          <Text className="text-lg font-semibold">{event.title}</Text>
          <Text className="text-grey-msg text-md mb-2">
            {event.description}
          </Text>
          <Text className="text-grey-msg mt-1">
            📅 {event.start_date} {event.end_date ? `— ${event.end_date}` : ""}
          </Text>
          <Text className="text-grey-msg">
            📍 {event.format === "online" ? "Онлайн" : event.address}
          </Text>
          <Text className="text-grey-msg">
            💰 Мін. внесок:{" "}
            {event.minimum_donation === 0
              ? "Безкоштовно"
              : `${event.minimum_donation} грн`}
          </Text>
        </View>
      </TouchableOpacity>

      <TouchableOpacity
        onPress={navigateToEvent}
        className="bg-primary rounded-full mt-6 py-3 items-center shadow-md"
      >
        <Text className="text-white text-lg font-semibold">
          Зареєструватись 📝
        </Text>
      </TouchableOpacity>
    </View>
  );
}
