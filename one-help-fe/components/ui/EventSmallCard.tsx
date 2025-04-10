import React from "react";
import { View, Text, TouchableOpacity } from "react-native";
import { useRouter } from "expo-router";

interface EventCardProps {
  event: {
    id: string;
    date: string;
    month: string;
    title: string;
    location: string;
    color: string;
  };
}

export default function EventSmallCard({ event }: EventCardProps) {
  const router = useRouter();

  const handlePress = () => {
    router.push(`/events/${event.id}`);
  };

  return (
    <View className="bg-accent p-4 rounded-2xl w-64 mr-4">
      <View className="flex-row items-center">
        <View className="bg-purple-bgr p-2 rounded-lg mr-3">
          <Text className="text-purple-msg text-sm">{event.month}</Text>
          <Text className="text-purple-msg text-2xl">{event.date}</Text>
        </View>
        <View>
          <Text className="text-black text-lg font-semibold">
            {event.title}
          </Text>
          <Text className="text-grey-msg">{event.location}</Text>
        </View>
      </View>

      <TouchableOpacity
        className="bg-primary py-2 rounded-xl mt-3"
        onPress={handlePress}
      >
        <Text className="text-white text-center font-semibold">
          Зареєструватись 📝
        </Text>
      </TouchableOpacity>
    </View>
  );
}
