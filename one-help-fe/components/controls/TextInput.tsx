import React from "react";
import {
  Text,
  View,
  TextInput as Input,
  TextInputProps as InputProps,
} from "react-native";

export interface TextInputProps extends InputProps {
  label?: string;
  required?: string;
  error?: boolean;
  errorMessage?: string;
}

export function TextInput({
  label,
  required,
  className,
  error = false, // Default value for error
  errorMessage,
  ...props
}: TextInputProps) {
  return (
    <View>
      {label && (
        <Text className="text-gray-msg font-medium">
          {label}{" "}
          {required && (
            <Text className="text-primary font-medium items-right">
              {required}
            </Text>
          )}
        </Text>
      )}
      <Input
        className={`border ${
          error ? "border-red-500" : "border-gray-300"
        } rounded-lg p-3 mt-1 ${className}`}
        {...props}
      />
      {error && errorMessage && (
        <Text className="text-red-500 text-sm mt-1">{errorMessage}</Text>
      )}
    </View>
  );
}
