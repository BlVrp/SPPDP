import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";

interface RaffleSmallCard {
  raffle: {
    image: string;
    title: string;
    description: string;
    minDonation: number;
  };
}

export default function RaffleSmallCard({ raffle }: RaffleSmallCard) {
  return (
    <View className="bg-accent p-4 rounded-2xl w-72 mr-4 mb-4 flex flex-col justify-between">
      <View className="flex-row items-center">
        <Image
          source={{ uri: raffle.image }}
          className="w-16 h-16 rounded-lg mr-4"
        />
        <View className="flex-1">
          <Text className="text-black text-lg font-semibold">
            {raffle.title}
          </Text>
          <Text className="text-grey-msg text-sm">{raffle.description}</Text>
          <Text className="text-green-600 text-sm font-semibold">
            Донат від: {raffle.minDonation} ₴
          </Text>
        </View>
      </View>

      <TouchableOpacity className="bg-primary py-2 rounded-xl mt-3">
        <Text className="text-white text-center font-semibold">
          Взяти участь 🎟️
        </Text>
      </TouchableOpacity>
    </View>
  );
}
