import { Link } from "expo-router";
import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";
import { ProgressBar } from "react-native-paper";

interface FundraiserCardProps {
  fundraiser: {
    id: string;
    title: string;
    description: string;
    image: string;
    goal: number;
    raised: number;
  };
}

export default function FundraiserCard({ fundraiser }: FundraiserCardProps) {
  const progress = fundraiser.raised / fundraiser.goal;

  return (
    <Link
      href={{
        pathname: "/fundraises/[id]",
        params: { id: `${fundraiser.id}` },
      }}
    >
      <View className="bg-white rounded-lg p-4 mb-4">
        <Image
          source={{ uri: fundraiser.image }}
          className="w-full h-32 rounded-lg mb-3"
          resizeMode="cover"
        />

        <Text className="text-lg font-semibold">{fundraiser.title}</Text>
        <Text className="text-grey-msg">{fundraiser.description}</Text>

        <Text className="text-grey-msg mt-2">
          {fundraiser.raised.toLocaleString()} /{" "}
          {fundraiser.goal.toLocaleString()}
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
