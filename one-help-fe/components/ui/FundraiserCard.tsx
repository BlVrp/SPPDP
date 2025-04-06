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
    image?: string; // optional
  };
}

export default function FundraiserCard({ fundraiser }: FundraiserCardProps) {
  const {
    id,
    title,
    description,
    targetAmount,
    filledAmount,
    image,
  } = fundraiser;

  const progress = Math.min(filledAmount / targetAmount || 0, 1);

  return (
    <Link
      href={{
        pathname: "/fundraises/[id]",
        params: { id },
      }}
    >
      <View className="w-full bg-white rounded-lg p-4 mb-4 shadow">
      
      <Image
  source={fundraiser.image?.length ? { uri: fundraiser.image } : defaultImage}
  className="w-full h-32 rounded-lg mb-3 self-center"
  resizeMode="cover"
/>

        <Text className="text-lg font-semibold text-black">{title}</Text>
        <Text className="text-gray-600 mt-1">{description}</Text>

        <Text className="text-gray-600 mt-2">
          {filledAmount.toLocaleString()} / {targetAmount.toLocaleString()} грн
        </Text>

        <ProgressBar
          progress={progress}
          color="#2563EB"
          className="h-2 rounded-xl mt-1"
        />

        <TouchableOpacity className="bg-primary rounded-lg p-2 mt-3 items-center">
          <Text className="text-white font-semibold">Донат 🍩</Text>
        </TouchableOpacity>
      </View>
    </Link>
  );
}
