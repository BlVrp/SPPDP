import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";
import { useRouter } from "expo-router";
import defaultImage from "@/assets/images/field.jpeg";

interface RaffleSmallCardProps {
  raffle: {
    id: string;
    title: string;
    description: string;
    minimumDonation: number;
    gifts: {
      imageUrl?: string;
    }[];
  };
}

export default function RaffleSmallCard({ raffle }: { raffle: RaffleSmallCardProps["raffle"] }) {
  const router = useRouter();
  const imageUrl = raffle.gifts?.[0]?.imageUrl || null;

  return (
    <View className="bg-accent p-4 rounded-2xl w-72 mr-4 mb-4 flex flex-col justify-between">
      <View className="flex-row items-center">
        <Image
          source={imageUrl ? { uri: imageUrl } : defaultImage}
          className="w-16 h-16 rounded-lg mr-4"
          resizeMode="cover"
        />
        <View className="flex-1">
          <Text className="text-black text-lg font-semibold">
            {raffle.title}
          </Text>
          <Text className="text-grey-msg text-sm">{raffle.description}</Text>
          <Text className="text-green-600 text-sm font-semibold">
            Донат від: {raffle.minimumDonation} ₴
          </Text>
        </View>
      </View>

      <TouchableOpacity
        className="bg-primary py-2 rounded-xl mt-3"
        onPress={() => router.push(`/raffles/${raffle.id}`)}
      >
        <Text className="text-white text-center font-semibold">
          Взяти участь 🎟️
        </Text>
      </TouchableOpacity>
    </View>
  );
}
