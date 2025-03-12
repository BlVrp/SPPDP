import {
View,
Text,
TouchableOpacity,
Modal,
TouchableWithoutFeedback,
Animated,
PanResponder,
} from "react-native";
import { useEffect, useRef } from "react";
import Icon from "react-native-vector-icons/Feather";
import { StatusBar } from "react-native";

interface SidebarProps {
isOpen: boolean;
onClose: () => void;
}

export default function SideBar({ isOpen, onClose }: SidebarProps) {
const translateX = useRef(new Animated.Value(-300)).current; // Початкове положення за екраном

useEffect(() => {
    Animated.timing(translateX, {
    toValue: isOpen ? 0 : -300,
    duration: 300,
    useNativeDriver: true,
    }).start();
}, [isOpen]);

const panResponder = useRef(
    PanResponder.create({
    onMoveShouldSetPanResponder: (_, gestureState) =>
        Math.abs(gestureState.dx) > 10, // Почати свайп при горизонтальному русі
    onPanResponderMove: (_, gestureState) => {
        if (gestureState.dx < 0) {
        translateX.setValue(gestureState.dx);
        }
    },
    onPanResponderRelease: (_, gestureState) => {
        if (gestureState.dx < -100) {
        onClose();
        } else {
        Animated.timing(translateX, {
            toValue: 0,
            duration: 200,
            useNativeDriver: true,
        }).start();
        }
    },
    })
).current;

return (
    <Modal visible={isOpen} animationType="fade" transparent>
    <TouchableWithoutFeedback onPress={onClose}>
        <View className="flex-1 bg-black bg-opacity-50">
        {/* Sidebar Container */}
        <Animated.View
            {...panResponder.panHandlers}
            style={{
            transform: [{ translateX }],
            width: "75%",
            height: "100%",
            backgroundColor: "white",
            paddingTop: StatusBar.currentHeight || 55,
            paddingHorizontal: 22,
            shadowColor: "#000",
            shadowOpacity: 0.2,
            shadowRadius: 4,
            }}
        >
            {/* Close Button */}
            <TouchableOpacity onPress={onClose} className="self-end mb-4">
            <Icon name="x" size={28} color="black" />
            </TouchableOpacity>

            {/* Sidebar Menu */}
            <TouchableOpacity onPress={() => console.log("Збори натиснуто")} className="py-2">
            <Text className="text-2xl font-semibold">💰 Збори</Text>
            </TouchableOpacity>

            <TouchableOpacity onPress={() => console.log("Події натиснуто")} className="py-2">
            <Text className="text-2xl pt-4 font-semibold">🪩 Події</Text>
            </TouchableOpacity>

            <TouchableOpacity onPress={() => console.log("Розіграші натиснуто")} className="py-2">
            <Text className="text-2xl pt-4 font-semibold">🎟 Розіграші</Text>
            </TouchableOpacity>

            <TouchableOpacity onPress={() => console.log("Налаштування натиснуто")} className="py-2">
            <Text className="text-2xl pt-4 font-semibold">️⚙ Налаштування</Text>
            </TouchableOpacity>
        </Animated.View>
        </View>
    </TouchableWithoutFeedback>
    </Modal>
);
}