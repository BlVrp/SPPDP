import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";
import { ProgressBar } from "react-native-paper";

interface FundraiseSmallCardProps {
  fundraiser: {
    title: string;
    description: string;
    image: string;
    goal: number;
    raised: number;
  };
}

export default function FundraiseSmallCard({
  fundraiser,
}: FundraiseSmallCardProps) {
  const progress = fundraiser.raised / fundraiser.goal;

  return (
    <View className="bg-white p-4 rounded-2xl mb-4">
      <View className="flex-row items-center">
        <Image
          source={{ uri: fundraiser.image }}
          className="w-16 h-16 rounded-lg mr-4"
        />
        <View className="flex-1">
          <Text className="text-lg font-bold text-black">
            {fundraiser.title}
          </Text>
          <Text className="text-grey-msg text-md">
            {fundraiser.description}
          </Text>
        </View>
      </View>

      <Text className="text-grey-msg mt-1 text-center">
        {fundraiser.raised.toLocaleString()} /{" "}
        {fundraiser.goal.toLocaleString()}
      </Text>
      <ProgressBar
        progress={progress}
        color="#2563EB"
        className="mt-2 h-3 rounded-lg"
      />

      <TouchableOpacity className="bg-blue-600 py-2 rounded-xl mt-2">
        <Text className="text-white text-center font-semibold text-md">
          Донат 🍩
        </Text>
      </TouchableOpacity>
    </View>
  );
}
