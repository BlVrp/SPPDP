import React from "react";
import { View, Text, TouchableOpacity, Image } from "react-native";
import { useRouter, Link } from "expo-router";

export default function RaffleCard({ raffle }: { raffle: any }) {
  const router = useRouter();

  return (
    <View className="bg-accent p-4 rounded-lg mb-4">
      <Link
        href={{
          pathname: "/raffles/[id]",
          params: { id: `${raffle.raffle_id}` },
        }}
      >
        <View className="w-full">
          {raffle.image && (
            <Image
              source={{ uri: raffle.image }}
              className="w-full h-32 rounded-lg mb-3"
              resizeMode="cover"
            />
          )}
          <Text className="text-lg font-semibold">{raffle.title}</Text>
          <Text className="text-grey-msg text-md mb-2">
            {raffle.description}
          </Text>
          <Text className="text-grey-msg mt-1">
            📅 {raffle.start_date}{" "}
            {raffle.end_date ? `— ${raffle.end_date}` : ""}
          </Text>
          <Text className="text-grey-msg">
            💰 Мін. внесок:{" "}
            {raffle.minimum_donation === 0
              ? "Безкоштовно"
              : `${raffle.minimum_donation} грн`}
          </Text>

          <TouchableOpacity
            onPress={() => console.log("Участь у розіграші")}
            className="bg-primary rounded-md p-2 mt-4 items-center"
          >
            <Text className="text-white text-lg font-semibold">
              Взяти участь 🎟️
            </Text>
          </TouchableOpacity>
        </View>
      </Link>
    </View>
  );
}
