import React from "react";
import { View, Text, TouchableOpacity, ScrollView } from "react-native";
import { useRouter } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { TextInput } from "@/components/controls";

export default function CreateFundraise() {
  const router = useRouter();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm();

  const onSubmit = (data: any) => {
    console.log("Новий збір:", data);
    router.back(); // Повертаємось на головний екран після створення
  };

  return (
    <ScrollView className="flex-1 bg-white p-6">
      <Text className="text-2xl font-bold text-gray-msg mb-4">
        Створити збір
      </Text>
      <View className="flex flex-col gap-4">
        <Controller
          control={control}
          name="title"
          rules={{ required: "Назва збору є обов'язковою" }}
          render={({
            field: { onChange, value },
          }: {
            field: { onChange: (value: string) => void; value: string };
          }) => (
            <TextInput
              label="Назва збору"
              required="*"
              placeholder="Дрон для бригади"
              onChangeText={onChange}
              value={value}
              error={!!errors.title}
              errorMessage={errors.title?.message?.toString()}
            />
          )}
        />

        <Controller
          control={control}
          name="description"
          rules={{ required: "Опис є обов'язковим" }}
          render={({
            field: { onChange, value },
          }: {
            field: { onChange: (value: string) => void; value: string };
          }) => (
            <TextInput
              label="Опис"
              placeholder="Детально опишіть, на що збір"
              required="*"
              multiline
              onChangeText={onChange}
              value={value}
              error={!!errors.description}
              errorMessage={errors.description?.message?.toString()}
            />
          )}
        />

        <Controller
          control={control}
          name="goal"
          rules={{ required: "Необхідна сума є обов'язковою" }}
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Необхідна сума (грн)"
              required="*"
              placeholder="50000"
              keyboardType="numeric"
              onChangeText={onChange}
              value={value}
              error={!!errors.goal}
              errorMessage={errors.goal?.message?.toString()}
            />
          )}
        />

        <Controller
          control={control}
          name="image"
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Посилання на фото"
              placeholder="https://example.com/image.jpg"
              onChangeText={onChange}
              value={value}
              error={!!errors.image}
              errorMessage={errors.image?.message?.toString()}
            />
          )}
        />

        <Controller
          control={control}
          name="endDate"
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Дата завершення"
              placeholder="27.05.2025"
              keyboardType="numeric"
              onChangeText={onChange}
              value={value}
              error={!!errors.endDate}
              errorMessage={errors.endDate?.message?.toString()}
            />
          )}
        />
      </View>

      <TouchableOpacity
        onPress={handleSubmit(onSubmit)}
        className="bg-primary rounded-lg p-4 mt-6 items-center"
      >
        <Text className="text-white text-lg font-semibold">
          Створити збір ✅
        </Text>
      </TouchableOpacity>
    </ScrollView>
  );
}
