import React from "react";
import { View, Text, TouchableOpacity } from "react-native";

interface EventCardProps {
  event: {
    date: string;
    month: string;
    title: string;
    location: string;
    color: string;
  };
}

export default function EventSmallCard({ event }: EventCardProps) {
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

      <TouchableOpacity className="bg-primary py-2 rounded-xl mt-3">
        <Text className="text-white text-center font-semibold">
          Зареєструватись 📝
        </Text>
      </TouchableOpacity>
    </View>
  );
}
