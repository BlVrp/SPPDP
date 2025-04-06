import React from "react";
import { View, Text, Image, TouchableOpacity } from "react-native";
import { ProgressBar } from "react-native-paper";
import defaultImage from "@/assets/images/field.jpeg";
import { useRouter } from "expo-router";

interface FundraiseSmallCardProps {
  fundraiser: {
    id?: string;
    title: string;
    description: string;
    image?: string;
    targetAmount?: number;
    filledAmount?: number;
  };
}

export default function FundraiseSmallCard({
  fundraiser,
}: FundraiseSmallCardProps) {
  const raised = fundraiser.filledAmount ?? 0;
  const goal = fundraiser.targetAmount ?? 1;
  const progress = raised / goal;
  const router = useRouter();

  return (
    <View className="bg-white p-4 rounded-2xl mb-4">
      <View className="flex-row items-center">
        <Image
          source={
            fundraiser.image?.length ? { uri: fundraiser.image } : defaultImage
          }
          style={{
            width: 54,
            height: 54,
            borderRadius: 8,
            marginRight: 12,
          }}
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
        {raised.toLocaleString()} / {goal.toLocaleString()}
      </Text>

      <ProgressBar
        progress={progress}
        color="#2563EB"
        className="mt-2 h-3 rounded-lg"
      />

<TouchableOpacity
        className="bg-blue-600 py-2 rounded-xl mt-2"
        onPress={() => router.push(`/fundraises/${fundraiser.id}`)}
      >
        <Text className="text-white text-center font-semibold text-md">
          Донат 🍩
        </Text>
      </TouchableOpacity>
    </View>
  );
}
