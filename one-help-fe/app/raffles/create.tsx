import React, { useEffect, useState } from "react";
import { View, Text, TouchableOpacity, ScrollView } from "react-native";
import { Dropdown } from "react-native-element-dropdown";
import { useRouter } from "expo-router";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { TextInput, Select } from "@/components/controls";

export default function CreateRaffle() {
  const router = useRouter();

  // Sample fundraisers list (replace this with data fetching logic)
  const [fundraisers, setFundraisers] = useState([
    { id: "1", name: "Fundraiser 1" },
    { id: "2", name: "Fundraiser 2" },
    { id: "3", name: "Fundraiser 3" },
  ]);

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

  const onSubmit = (data: any) => {
    console.log("Новий розіграш:", data);
    router.back();
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
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Опис"
              required="*"
              placeholder="Опишіть деталі розіграшу"
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
          name="minimum_donation"
          rules={{ required: "Мінімальний донат є обов'язковим" }}
          render={({ field: { onChange, value } }) => (
            <TextInput
              label="Мінімальний донат (₴)"
              required="*"
              placeholder="50"
              keyboardType="numeric"
              onChangeText={onChange}
              value={value}
              error={!!errors.minimum_donation}
              errorMessage={errors.minimum_donation?.message?.toString()}
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
              labelField="name"
              valueField="id"
              label="Збір"
              placeholder="Оберіть збір"
              searchPlaceholder="Search..."
              value={value}
              onChange={onChange}
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
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.title}
                  errorMessage={errors.gifts?.[
                    index
                  ]?.title?.message?.toString()}
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
                  multiline
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.description}
                  errorMessage={errors.gifts?.[
                    index
                  ]?.description?.message?.toString()}
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
                  onChangeText={onChange}
                  value={value}
                  error={!!errors.gifts?.[index]?.image}
                  errorMessage={errors.gifts?.[
                    index
                  ]?.image?.message?.toString()}
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
          className="bg-primary rounded-lg p-4 items-center"
        >
          <Text className="text-white text-lg font-semibold">
            Створити розіграш ✅
          </Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}
