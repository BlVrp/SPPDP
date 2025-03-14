import React, { useState } from "react";
import {
  FlatListProps,
  ImageStyle,
  StyleProp,
  StyleSheet,
  Text,
  TextProps,
  TextStyle,
  View,
  ViewStyle,
} from "react-native";
import { Dropdown, IDropdownRef } from "react-native-element-dropdown";

export interface DropdownProps<T> {
  label?: string;
  required?: boolean;
  ref?:
    | React.RefObject<IDropdownRef>
    | React.MutableRefObject<IDropdownRef>
    | null
    | undefined;
  testID?: string;
  itemTestIDField?: string;
  style?: StyleProp<ViewStyle>;
  containerStyle?: StyleProp<ViewStyle>;
  placeholderStyle?: StyleProp<TextStyle>;
  selectedTextStyle?: StyleProp<TextStyle>;
  selectedTextProps?: TextProps;
  itemContainerStyle?: StyleProp<ViewStyle>;
  itemTextStyle?: StyleProp<TextStyle>;
  inputSearchStyle?: StyleProp<TextStyle>;
  iconStyle?: StyleProp<ImageStyle>;
  maxHeight?: number;
  minHeight?: number;
  fontFamily?: string;
  iconColor?: string;
  activeColor?: string;
  data: T[];
  value?: T | string | null | undefined;
  placeholder?: string;
  labelField: keyof T;
  valueField: keyof T;
  searchField?: keyof T;
  error?: boolean;
  errorMessage?: string;
  search?: boolean;
  searchPlaceholder?: string;
  searchPlaceholderTextColor?: string;
  disable?: boolean;
  autoScroll?: boolean;
  showsVerticalScrollIndicator?: boolean;
  dropdownPosition?: "auto" | "top" | "bottom";
  flatListProps?: Omit<FlatListProps<any>, "renderItem" | "data">;
  keyboardAvoiding?: boolean;
  backgroundColor?: string;
  confirmSelectItem?: boolean;
  accessibilityLabel?: string;
  itemAccessibilityLabelField?: string;
  inverted?: boolean;
  mode?: "default" | "modal" | "auto";
  closeModalWhenSelectedItem?: boolean;
  excludeItems?: T[];
  excludeSearchItems?: T[];
  onChange: (item: T) => void;
  renderLeftIcon?: (visible?: boolean) => React.ReactElement | null;
  renderRightIcon?: (visible?: boolean) => React.ReactElement | null;
  renderItem?: (item: T, selected?: boolean) => React.ReactElement | null;
  renderInputSearch?: (
    onSearch: (text: string) => void
  ) => React.ReactElement | null;
  onFocus?: () => void;
  onBlur?: () => void;
  searchQuery?: (keyword: string, labelValue: string) => boolean;
  onChangeText?: (search: string) => void;
  onConfirmSelectItem?: (item: T) => void;
}

export const Select = ({
  value,
  label,
  required,
  error = false,
  errorMessage,
  ...props
}: DropdownProps<any>) => {
  const [isFocus, setIsFocus] = useState(false);

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
      <View style={styles.container}>
        <Dropdown
          style={[
            styles.dropdown,
            isFocus && { borderColor: "blue" },
            error && { borderColor: "red" },
          ]}
          placeholderStyle={styles.placeholderStyle}
          selectedTextStyle={styles.selectedTextStyle}
          inputSearchStyle={styles.inputSearchStyle}
          iconStyle={styles.iconStyle}
          placeholder={!isFocus ? "Select item" : "..."}
          searchPlaceholder="Search..."
          value={value}
          onFocus={() => setIsFocus(true)}
          onBlur={() => setIsFocus(false)}
          {...props}
        />
        {error && errorMessage && (
          <Text className="text-red-500 text-sm mt-1">{errorMessage}</Text>
        )}
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: "white",
  },
  dropdown: {
    height: 50,
    borderColor: "gray",
    borderWidth: 0.5,
    borderRadius: 8,
    paddingHorizontal: 8,
  },
  icon: {
    marginRight: 5,
  },
  label: {
    position: "absolute",
    backgroundColor: "white",
    left: 22,
    top: 8,
    zIndex: 999,
    paddingHorizontal: 8,
    fontSize: 14,
  },
  placeholderStyle: {
    fontSize: 16,
  },
  selectedTextStyle: {
    fontSize: 16,
  },
  iconStyle: {
    width: 20,
    height: 20,
  },
  inputSearchStyle: {
    height: 40,
    fontSize: 16,
  },
});
