import React, { useState } from "react";
import { View, Text, Image, TouchableOpacity, ScrollView } from "react-native";
import { useLocalSearchParams } from "expo-router";

const TEMP_RAFFLE = {
  raffle_id: "1",
  title: 'Збір "Ліжка для війська"',
  description:
    "Кожен донат бере участь у розіграші 2 принтів картини «Свята земля» з автографами KALUSH ORCHESTRA",
  minimum_donation: 10,
  participation_status: false,
  gifts: [
    {
      id: "1",
      title: "Playstation",
      image:
        "https://kolo-django-prod.s3.amazonaws.com/campaigns/2/1692702564/ps_5_1200%C3%91630.webp",
    },
    {
      id: "2",
      title: "Сережки з гербом",
      description: "",
      image:
        "https://soundmagcdn.fra1.cdn.digitaloceanspaces.com/news/688/uG5AEqgmTpT61F7n.webp",
    },
  ],
  fundraiser: {
    target: 2500000,
    current: 1000000,
    end_date: "2025-11-30",
  },
};

export default function RaffleDetail() {
  const { id } = useLocalSearchParams();
  const [currentGiftIndex, setCurrentGiftIndex] = useState(0);
  const raffle = TEMP_RAFFLE;

  const nextGift = () => {
    setCurrentGiftIndex((prev) =>
      prev + 1 < raffle.gifts.length ? prev + 1 : 0
    );
  };

  const prevGift = () => {
    setCurrentGiftIndex((prev) =>
      prev - 1 >= 0 ? prev - 1 : raffle.gifts.length - 1
    );
  };

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-4 rounded-2xl">
        <View className="relative items-center">
          <TouchableOpacity
            onPress={prevGift}
            className="absolute left-0 top-1/2 -translate-y-1/2 px-3 py-2"
          >
            <Text className="text-2xl">⬅</Text>
          </TouchableOpacity>

          <Image
            source={{ uri: raffle.gifts[currentGiftIndex].image }}
            className="w-3/4 h-72 rounded-lg"
            resizeMode="cover"
          />
          <Text className="text-lg text-grey-msg text-center mt-1">
            {raffle.gifts[currentGiftIndex].title}
          </Text>

          <TouchableOpacity
            onPress={nextGift}
            className="absolute right-0 top-1/2 -translate-y-1/2 px-3 py-2"
          >
            <Text className="text-2xl">➡</Text>
          </TouchableOpacity>
        </View>

        <View className="mt-4 flex-row justify-center items-center">
          <TouchableOpacity className="bg-primary rounded-md px-4 w-2/3 py-2">
            <Text className="text-white text-lg font-semibold text-center">
              Донат від {raffle.minimum_donation} ₴ 💰
            </Text>
          </TouchableOpacity>

          <Text className="text-4xl ml-3">
            {raffle.participation_status ? "✅" : "☑️"}
          </Text>
        </View>

        <Text className="text-2xl font-bold text-black text-center mt-4">
          {raffle.title}
        </Text>
        <View className="mb-3 ml-4">
          <Text className="text-grey-msg text-md mt-4">
            {raffle.description}
          </Text>

          <Text className="text-grey-msg mt-2">
            📆 Донати приймаються до {raffle.fundraiser.end_date}
          </Text>
        </View>

        <View className="mt-4">
          <Text className="text-grey-msg text-center">
            {raffle.fundraiser.current.toLocaleString()} /{" "}
            {raffle.fundraiser.target.toLocaleString()}
          </Text>
          <View className="h-2 w-full bg-gray-200 rounded-full mt-1">
            <View
              className="h-2 bg-primary rounded-full"
              style={{
                width: `${
                  (raffle.fundraiser.current / raffle.fundraiser.target) * 100
                }%`,
              }}
            />
          </View>
        </View>
      </View>
    </ScrollView>
  );
}
