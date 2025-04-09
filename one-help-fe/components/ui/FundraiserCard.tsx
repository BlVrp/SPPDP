import { Link } from "expo-router";
import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";
import { ProgressBar } from "react-native-paper";
import defaultImage from "@/assets/images/field.jpeg";

interface FundraiserCardProps {
  fundraiser: {
    id: string;
    title: string;
    description: string;
    targetAmount: number;
    filledAmount: number;
    imageUrl: string;
    endDate: string;
  };
}

export default function FundraiserCard({ fundraiser }: FundraiserCardProps) {
  const { id, title, description, targetAmount, filledAmount, imageUrl } =
    fundraiser;

  const progress = Math.min(filledAmount / targetAmount || 0, 1);

  return (
    <Link
      href={{
        pathname: "/fundraises/[id]",
        params: { id },
      }}
    >
      <View className="w-full bg-accent p-4 rounded-lg mb-4">
        <Image
          source={imageUrl?.length ? { uri: imageUrl } : defaultImage}
          className="w-full h-32 rounded-lg mb-3 self-center"
          resizeMode="cover"
        />

        <Text className="text-lg font-semibold text-black">{title}</Text>
        <Text className="text-gray-600 mt-1">{description}</Text>

        <View className="mt-2">
          <ProgressBar
            progress={progress}
            color="#2563EB"
            style={{ height: 10, borderRadius: 8 }}
          />
          <Text className="text-gray-600 mt-1 text-sm text-center">
            {filledAmount.toLocaleString()} / {targetAmount.toLocaleString()}{" "}
            грн
          </Text>
        </View>

        <TouchableOpacity className="bg-primary rounded-full mt-2 py-3 items-center shadow-md">
          <Text className="text-white font-semibold">Донат 🍩</Text>
        </TouchableOpacity>
      </View>
    </Link>
  );
}
