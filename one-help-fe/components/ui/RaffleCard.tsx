import React from "react";
import { View, Text, TouchableOpacity, Image } from "react-native";
import { useRouter, Link } from "expo-router";

function formatDate(isoDateString: string) {
  const dateObj = new Date(isoDateString);
  return dateObj.toLocaleString("uk-UA", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

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
            📅 {formatDate(raffle.start_date)} 
            {raffle.end_date ? ` — ${formatDate(raffle.end_date)}` : ""}
          </Text>

          <Text className="text-grey-msg">
            💰 Мін. внесок:{" "}
            {raffle.minimum_donation === 0
              ? "Безкоштовно"
              : `${raffle.minimum_donation} грн`}
          </Text>

          <TouchableOpacity
            className="bg-primary rounded-md p-1 mt-4 items-center"
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
