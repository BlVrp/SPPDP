import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function RaffleDetail() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const [raffle, setRaffle] = useState<any>(null);
  const [currentGiftIndex, setCurrentGiftIndex] = useState(0);
  const [loading, setLoading] = useState(true);

  const [participants, setParticipants] = useState<any[]>([]);
  const [showParticipants, setShowParticipants] = useState(false);

  const fetchRaffle = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const res = await fetch(`http://localhost:8080/api/v0/raffles/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Не вдалося отримати розіграш");
      }

      const result = await res.json();
      setRaffle(result);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchParticipants = async () => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) return;

      const res = await fetch(
        `http://localhost:8080/api/v0/users/raffle-participants/${id}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!res.ok) {
        const error = await res.json();
        throw new Error(error.message || "Не вдалося отримати учасників");
      }

      const data = await res.json();
      setParticipants(data);
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    }
  };

  const toggleParticipants = async () => {
    if (!showParticipants && participants.length === 0) {
      await fetchParticipants();
    }
    setShowParticipants(!showParticipants);
  };

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

  useEffect(() => {
    if (id) fetchRaffle();
  }, [id]);

  if (loading) {
    return (
      <View className="flex-1 justify-center items-center bg-white">
        <ActivityIndicator size="large" color="#2563EB" />
      </View>
    );
  }

  if (!raffle) {
    return (
      <View className="flex-1 justify-center items-center bg-white p-6">
        <Text className="text-lg text-gray-msg">Розіграш не знайдено 😕</Text>
      </View>
    );
  }

  return (
    <ScrollView className="flex-1 bg-white p-4">
      <View className="bg-accent p-5 rounded-3xl shadow-sm">
        {/* Gift Navigation */}
        {raffle.gifts.length > 1 && (
          <View className="flex-row justify-center items-center space-x-8 mb-2">
            <TouchableOpacity onPress={prevGift}>
              <Text className="text-3xl">⬅</Text>
            </TouchableOpacity>
            <Text className="text-sm text-gray-500">
              {currentGiftIndex + 1} / {raffle.gifts.length}
            </Text>
            <TouchableOpacity onPress={nextGift}>
              <Text className="text-3xl">➡</Text>
            </TouchableOpacity>
          </View>
        )}

        {/* Gift Image and Title */}
        {raffle.gifts.length > 0 && (
          <View className="items-center">
            <Image
              source={{ uri: raffle.gifts[currentGiftIndex].imageUrl }}
              className="w-full h-64 rounded-xl mb-2"
              resizeMode="cover"
            />
            <Text className="text-lg font-medium text-gray-700 text-center mb-2">
              🎁 {raffle.gifts[currentGiftIndex].title}
            </Text>
          </View>
        )}

        {/* Donate Button */}
        <TouchableOpacity
          className="bg-primary rounded-full py-3 mt-4 items-center"
          onPress={() => router.push(`/fundraises/${raffle.fundraiseId}`)}
        >
          <Text className="text-white text-base font-semibold">
            Донат від {raffle.minimumDonation} ₴ 💰
          </Text>
        </TouchableOpacity>

        {/* Raffle Title */}
        <Text className="text-2xl font-bold text-black text-center mt-6">
          {raffle.title}
        </Text>

        {/* Raffle Description */}
        <Text className="text-gray-600 text-base text-center mt-3">
          {raffle.description}
        </Text>

        {/* Toggle Participants */}
        <TouchableOpacity
          onPress={toggleParticipants}
          className=" mt-6 py-2 px-4 rounded-xl items-center"
        >
          <Text className="text-primary font-medium">
            {showParticipants ? "Сховати учасників" : "Показати учасників"}
          </Text>
        </TouchableOpacity>

        {/* Participant List */}
        {showParticipants && (
          <View className="mt-4 bg-white rounded-xl px-4 py-3 shadow-sm">
            <Text className="text-lg font-semibold mb-3 text-center">
              Учасники ({participants.length})
            </Text>
            {participants.length === 0 ? (
              <Text className="text-gray-500 text-center">
                Поки що немає учасників
              </Text>
            ) : (
              participants.map((p) => (
                <View key={p.id} className="flex-row items-center mb-3">
                  {p.imageUrl ? (
                    <Image
                      source={{ uri: p.imageUrl }}
                      className="w-10 h-10 rounded-full mr-3"
                      resizeMode="cover"
                    />
                  ) : (
                    <View className="w-10 h-10 bg-gray-300 rounded-full mr-3 items-center justify-center">
                      <Text className="text-white text-sm">
                        {p.firstName[0]}
                        {p.lastName[0]}
                      </Text>
                    </View>
                  )}
                  <View>
                    <Text className="text-black font-medium">
                      {p.firstName} {p.lastName}
                    </Text>
                    {p.city && (
                      <Text className="text-gray-500 text-sm">{p.city}</Text>
                    )}
                  </View>
                </View>
              ))
            )}
          </View>
        )}
      </View>
    </ScrollView>
  );
}
