import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  Alert,
} from "react-native";
import { useRouter } from "expo-router";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { TextInput, Select } from "@/components/controls";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function CreateRaffle() {
  const router = useRouter();
  const [fundraisers, setFundraisers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues: {
      title: "",
      description: "",
      minimum_donation: "",
      fundraiser_id: "",
      gifts: [{ title: "", description: "", image: "" }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "gifts",
  });

  useEffect(() => {
    const fetchMyFundraisers = async () => {
      try {
        const token = await AsyncStorage.getItem("token");
        if (!token) {
          return;
        }
        const response = await fetch("http://localhost:8080/api/v0/fundraises/my", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (!response.ok) {
          throw new Error("Не вдалося отримати ваші збори");
        }
        const result = await response.json();
        setFundraisers(result.data || []);
      } catch (error) {
        console.error("Failed to fetch user's fundraises:", error);
      } finally {
        setLoading(false);
      }
    };
    fetchMyFundraisers();
  }, []);

  const onSubmit = async (data: any) => {
    try {
      const token = await AsyncStorage.getItem("token");
      if (!token) {
        Alert.alert("Помилка", "Користувач не авторизований");
        return;
      }

      const response = await fetch("http://localhost:8080/api/v0/raffles/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          title: data.title,
          description: data.description,
          minimumDonation: Number(data.minimum_donation),
          fundraiseId: data.fundraiser_id,
          startDate: new Date().toISOString(),
          endDate: new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString(),
          gifts: data.gifts.map((gift: any) => ({
            title: gift.title,
            description: gift.description,
            imageUrl: gift.image,
          })),
        }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не вдалося створити розіграш");
      }

      Alert.alert("Успіх", "Розіграш створено успішно ✅");
      router.push("/raffles");
    } catch (err: any) {
      Alert.alert("Помилка", err.message);
    }
  };

  return (
    <ScrollView className="flex-1 bg-white p-6">
      <Text className="text-2xl font-bold text-gray-msg mb-4">
        Створити розіграш
      </Text>
      <View className="flex flex-col gap-4">
        <Controller
          control={control}
          name="title"
          rules={{ required: "Назва розіграшу є обов'язковою" }}
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Назва розіграшу"
              required="*"
              placeholder="Розіграш подарунків"
              placeholderTextColor="#9CA3AF"
              onChangeText={onChange}
              value={value}
              error={!!errors.title}
              errorMessage={errors.title?.message}
            />
          )}
        />

        <Controller
          control={control}
          name="description"
          rules={{ required: "Опис є обов'язковим" }}
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Опис"
              required="*"
              placeholder="Опишіть деталі розіграшу"
              placeholderTextColor="#9CA3AF"
              multiline
              onChangeText={onChange}
              value={value}
              error={!!errors.description}
              errorMessage={errors.description?.message}
            />
          )}
        />

        <Controller
          control={control}
          name="minimum_donation"
          rules={{ required: "Мінімальний донат є обов'язковим" }}
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Мінімальний донат (₴)"
              required="*"
              placeholder="50"
              placeholderTextColor="#9CA3AF"
              keyboardType="numeric"
              onChangeText={onChange}
              value={value}
              error={!!errors.minimum_donation}
              errorMessage={errors.minimum_donation?.message}
            />
          )}
        />

        <Controller
          control={control}
          name="fundraiser_id"
          rules={{ required: "Оберіть збір коштів" }}
          render={({ field: { onChange, value } }) => (
            <Select
              data={fundraisers}
              labelField="title"
              valueField="id"
              label="Збір"
              placeholder="Оберіть збір"
              onChange={(selected: any) => onChange(selected.id)}
              value={value}
              error={!!errors.fundraiser_id}
              errorMessage={errors?.fundraiser_id?.message}
            />
          )}
        />

        <Text className="text-lg font-semibold mb-2">🎁 Подарунки:</Text>
        {fields.map((gift, index) => (
          <View
            key={gift.id}
            className="border border-gray-300 p-3 rounded-md mb-3"
          >
            <Controller
              control={control}
              name={`gifts.${index}.title`}
              rules={{ required: "Назва подарунка є обов'язковою" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Назва подарунка"
                  required="*"
                  placeholder="Apple Watch"
                  placeholderTextColor="#9CA3AF"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.title}
                  errorMessage={errors.gifts?.[index]?.title?.message}
                />
              )}
            />

            <Controller
              control={control}
              name={`gifts.${index}.description`}
              rules={{ required: "Опис подарунка є обов'язковим" }}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Опис подарунка"
                  required="*"
                  placeholder="Опис подарунка..."
                  placeholderTextColor="#9CA3AF"
                  multiline
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.description}
                  errorMessage={errors.gifts?.[index]?.description?.message}
                />
              )}
            />

            <Controller
              control={control}
              name={`gifts.${index}.image`}
              render={({ field: { onChange, value } }) => (
                <TextInput
                  label="Посилання на зображення"
                  placeholder="https://example.com/gift.jpg"
                  placeholderTextColor="#9CA3AF"
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.image}
                  errorMessage={errors.gifts?.[index]?.image?.message}
                />
              )}
            />

            {fields.length > 1 && (
              <TouchableOpacity
                onPress={() => remove(index)}
                className="bg-red-500 rounded-md mt-2 p-2"
              >
                <Text className="text-white text-center">
                  Видалити подарунок
                </Text>
              </TouchableOpacity>
            )}
          </View>
        ))}

        <TouchableOpacity
          onPress={() => append({ title: "", description: "", image: "" })}
          className="bg-gray-300 rounded-md p-3 mb-3"
        >
          <Text className="text-black text-center">+ Додати ще подарунок</Text>
        </TouchableOpacity>

        <TouchableOpacity
          onPress={handleSubmit(onSubmit)}
          className="bg-primary rounded-lg p-1 items-center"
        >
          <Text className="text-white text-lg font-semibold">
            Створити розіграш ✅
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
