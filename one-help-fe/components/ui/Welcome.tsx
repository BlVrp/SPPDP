import React from "react";
import { View, Text } from "react-native";

interface WelcomeProps {
  user: {
    name: string;
    donated: number;
    events: number;
    raffles: number;
  };
}

const Welcome = ({ user }: WelcomeProps) => {
  return (
    <View className="bg-accent p-4 rounded-lg">
      <Text className="text-4xl text-black mt-6">Привіт, {user.name}!</Text>
      <Text className="text-2xl text-grey-msg mt-3 mb-4">
        Зміни починаються сьогодні
      </Text>

      <View className="flex-row justify-between mt-4 space-x-4 mb-8">
        <View className="bg-white p-4 rounded-xl items-left flex-1">
          <Text className="text-2xl text-grey-msg mb-1">🍩</Text>
          <Text className="text-xl text-grey-msg mb-2">Донатів</Text>
          <Text className="text-2xl text-black">₴ {user.donated}</Text>
        </View>
        <View className="bg-white p-4 rounded-xl mr-3 items-left flex-1">
          <Text className="text-2xl text-grey-msg mb-1">📆</Text>
          <Text className="text-xl text-grey-msg mb-2">Заходів</Text>
          <Text className="text-2xl text-black">{user.events}</Text>
        </View>
        <View className="bg-white p-4 rounded-xl mr-3 items-left flex-1">
          <Text className="text-2xl text-grey-msg mb-1">🎟</Text>
          <Text className="text-xl text-grey-msg mb-2">Розіграші</Text>
          <Text className="text-2xl text-black">{user.raffles}</Text>
        </View>
      </View>
    </View>
  );
};

export default Welcome;
