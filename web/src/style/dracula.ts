import { definePreset } from '@primeuix/themes'
import Aura from '@primeuix/themes/lara'

// Dracula theme colors
const dracula = {
    background: '#282a36',
    currentLine: '#44475a',
    selection: '#44475a',
    foreground: '#f8f8f2',
    comment: '#6272a4',
    red: '#ff5555',
    orange: '#ffb86c',
    yellow: '#f1fa8c',
    green: '#50fa7b',
    cyan: '#8be9fd',
    purple: '#bd93f9',
    pink: '#ff79c6',
    // Light theme specific colors
    light: {
        white: '#ffffff',
        lightGray: '#fafafa',
        borderGray: '#e5e7eb',
        mutedGray: '#d9dce3',
        darkText: '#1c1e28',
        mutedText: '#546293',
        lightPurple: '#eee9fd',
        focusPurple: '#e4d9fc',
        selectedPurple: '#ddd0fb',
        hoverPurple: '#c7a4fb',
        activePurple: '#a97cf5',
        disabledBg: '#f0f0f0'
    },
    // Dark theme specific colors
    dark: {
        darker: '#1e202b',
        darkest: '#232530',
        lighter: '#2e3040',
        border: '#373a4b',
        surface: '#51546a',
        muted: '#6b6e83',
        text: '#8a8fb1',
        lightText: '#b0b4d2',
        nearWhite: '#e6e8f5',
        hoverMuted: '#7b89c6'
    }
} as const

export default definePreset(Aura, {
    semantic: {
        primary: {
            50: '#f4eefe',
            100: '#eadbfe',
            200: '#d6b7fd',
            300: '#c393fb',
            400: '#b16ff9',
            500: dracula.purple,
            600: '#a671f2',
            700: '#8f55e3',
            800: '#7a3ecc',
            900: '#692fb3',
            950: '#4d1f8a'
        },
        colorScheme: {
            light: {
                surface: {
                    0: '#000000',
                    50: dracula.light.white,
                    100: dracula.light.lightGray,
                    200: dracula.foreground,
                    300: '#f0f1f5',
                    400: dracula.light.borderGray,
                    500: dracula.light.mutedGray,
                    600: '#cbd0db',
                    700: '#b8bfd0',
                    800: '#9ea6b9',
                    900: '#7f87a0',
                    950: '#5d6681'
                },
                primary: {
                    color: dracula.purple,
                    contrastColor: dracula.light.white,
                    hoverColor: dracula.light.hoverPurple,
                    activeColor: dracula.light.activePurple
                },
                highlight: {
                    background: dracula.light.lightPurple,
                    focusBackground: dracula.light.focusPurple,
                    color: dracula.background,
                    focusColor: dracula.background
                },
                focusRing: {
                    shadow: `0 0 0 2px ${dracula.purple}59`
                },
                mask: {
                    background: 'rgba(0,0,0,0.5)',
                    color: dracula.background
                },
                formField: {
                    background: dracula.light.white,
                    disabledBackground: dracula.light.disabledBg,
                    filledBackground: dracula.foreground,
                    filledHoverBackground: dracula.light.white,
                    filledFocusBackground: dracula.light.white,
                    borderColor: dracula.light.borderGray,
                    hoverBorderColor: dracula.light.mutedGray,
                    focusBorderColor: dracula.purple,
                    invalidBorderColor: dracula.red,
                    color: dracula.background,
                    disabledColor: dracula.comment,
                    placeholderColor: dracula.comment,
                    invalidPlaceholderColor: dracula.red,
                    floatLabelColor: dracula.comment,
                    floatLabelFocusColor: dracula.purple,
                    floatLabelActiveColor: dracula.comment,
                    floatLabelInvalidColor: dracula.red,
                    iconColor: dracula.background,
                    shadow: 'none'
                },
                text: {
                    color: dracula.background,
                    hoverColor: dracula.light.darkText,
                    mutedColor: dracula.comment,
                    hoverMutedColor: dracula.light.mutedText
                },
                content: {
                    background: dracula.light.white,
                    hoverBackground: dracula.foreground,
                    borderColor: dracula.light.borderGray,
                    color: dracula.background,
                    hoverColor: dracula.light.darkText
                },
                overlay: {
                    select: {
                        background: dracula.light.white,
                        borderColor: dracula.light.borderGray,
                        color: dracula.background
                    },
                    popover: {
                        background: dracula.light.white,
                        borderColor: dracula.light.borderGray,
                        color: dracula.background
                    },
                    modal: {
                        background: dracula.light.white,
                        borderColor: dracula.light.borderGray,
                        color: dracula.background
                    }
                },
                list: {
                    option: {
                        focusBackground: dracula.light.lightPurple,
                        selectedBackground: dracula.light.focusPurple,
                        selectedFocusBackground: dracula.light.selectedPurple,
                        color: dracula.background,
                        focusColor: dracula.background,
                        selectedColor: dracula.background,
                        selectedFocusColor: dracula.background,
                        icon: {
                            color: dracula.comment,
                            focusColor: dracula.background
                        }
                    },
                    optionGroup: {
                        background: dracula.light.white,
                        color: dracula.comment
                    }
                },
                navigation: {
                    item: {
                        focusBackground: dracula.light.lightPurple,
                        activeBackground: dracula.light.focusPurple,
                        color: dracula.background,
                        focusColor: dracula.background,
                        activeColor: dracula.background,
                        icon: {
                            color: dracula.comment,
                            focusColor: dracula.background,
                            activeColor: dracula.background
                        }
                    },
                    submenuLabel: {
                        background: dracula.light.white,
                        color: dracula.comment
                    },
                    submenuIcon: {
                        color: dracula.comment,
                        focusColor: dracula.background,
                        activeColor: dracula.background
                    }
                }
            },
            dark: {
                surface: {
                    0: dracula.light.white,
                    50: dracula.dark.darker,
                    100: dracula.dark.darkest,
                    200: dracula.background,
                    300: dracula.dark.lighter,
                    400: dracula.dark.border,
                    500: dracula.currentLine,
                    600: dracula.dark.surface,
                    700: dracula.dark.muted,
                    800: dracula.dark.text,
                    900: dracula.dark.lightText,
                    950: dracula.dark.nearWhite
                },
                primary: {
                    color: dracula.purple,
                    contrastColor: dracula.background,
                    hoverColor: dracula.light.hoverPurple,
                    activeColor: dracula.light.activePurple
                },
                highlight: {
                    background: dracula.currentLine,
                    focusBackground: dracula.dark.surface,
                    color: dracula.foreground,
                    focusColor: dracula.foreground
                },
                focusRing: {
                    shadow: `0 0 0 2px ${dracula.purple}66`
                },
                mask: {
                    background: `${dracula.background}cc`,
                    color: dracula.foreground
                },
                formField: {
                    background: dracula.background,
                    disabledBackground: dracula.dark.darkest,
                    filledBackground: dracula.dark.lighter,
                    filledHoverBackground: dracula.dark.border,
                    filledFocusBackground: dracula.dark.lighter,
                    borderColor: dracula.currentLine,
                    hoverBorderColor: dracula.dark.surface,
                    focusBorderColor: dracula.purple,
                    invalidBorderColor: dracula.red,
                    color: dracula.foreground,
                    disabledColor: dracula.comment,
                    placeholderColor: dracula.comment,
                    invalidPlaceholderColor: dracula.red,
                    floatLabelColor: dracula.comment,
                    floatLabelFocusColor: dracula.purple,
                    floatLabelActiveColor: dracula.comment,
                    floatLabelInvalidColor: dracula.red,
                    iconColor: dracula.foreground,
                    shadow: 'none'
                },
                text: {
                    color: dracula.foreground,
                    hoverColor: dracula.light.white,
                    mutedColor: dracula.comment,
                    hoverMutedColor: dracula.dark.hoverMuted
                },
                content: {
                    background: dracula.background,
                    hoverBackground: dracula.dark.lighter,
                    borderColor: dracula.currentLine,
                    color: dracula.foreground,
                    hoverColor: dracula.light.white
                },
                overlay: {
                    select: {
                        background: dracula.background,
                        borderColor: dracula.currentLine,
                        color: dracula.foreground
                    },
                    popover: {
                        background: dracula.background,
                        borderColor: dracula.currentLine,
                        color: dracula.foreground
                    },
                    modal: {
                        background: dracula.background,
                        borderColor: dracula.currentLine,
                        color: dracula.foreground
                    }
                },
                list: {
                    option: {
                        focusBackground: dracula.dark.border,
                        selectedBackground: dracula.currentLine,
                        selectedFocusBackground: dracula.dark.surface,
                        color: dracula.foreground,
                        focusColor: dracula.foreground,
                        selectedColor: dracula.foreground,
                        selectedFocusColor: dracula.foreground,
                        icon: {
                            color: dracula.foreground,
                            focusColor: dracula.light.white
                        }
                    },
                    optionGroup: {
                        background: dracula.background,
                        color: dracula.comment
                    }
                },
                navigation: {
                    item: {
                        focusBackground: dracula.dark.border,
                        activeBackground: dracula.currentLine,
                        color: dracula.foreground,
                        focusColor: dracula.light.white,
                        activeColor: dracula.light.white,
                        icon: {
                            color: dracula.foreground,
                            focusColor: dracula.light.white,
                            activeColor: dracula.light.white
                        }
                    },
                    submenuLabel: {
                        background: dracula.background,
                        color: dracula.comment
                    },
                    submenuIcon: {
                        color: dracula.comment,
                        focusColor: dracula.foreground,
                        activeColor: dracula.foreground
                    }
                }
            }
        }
    }
})